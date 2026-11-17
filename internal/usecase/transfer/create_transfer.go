package transfer

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/constant"
	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/entity"
	merchantrepo "github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/merchant"
	transferrepo "github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/transfer"
	"github.com/google/uuid"
)

func (u *useCase) CreateTransfer(ctx context.Context, req entity.CreateTransferRequest) (*entity.CreateTransferResponse, error) {
	idemCacheKey := fmt.Sprintf("idem:transfer:%d:%s", constant.DefaultMerchantID, req.IdempotencyKey)

	if cachedResp, hit, err := u.checkIdempotencyPreLock(ctx, idemCacheKey, req.IdempotencyKey); err != nil {
		return nil, err
	} else if hit {
		return cachedResp, nil
	}

	var created *entity.Transfer

	err := u.txManager.WithTx(ctx, func(txCtx context.Context) error {
		// Lock merchant row to serialize concurrent balance deductions.
		merchant, err := u.merchantRepo.GetByIDForUpdate(txCtx, constant.DefaultMerchantID)
		if err != nil {
			return err
		}

		if merchant.Balance < req.Amount {
			return merchantrepo.ErrInsufficientBalance
		}

		balanceAfter := merchant.Balance - req.Amount
		if err := u.merchantRepo.UpdateBalance(txCtx, merchant.ID, balanceAfter); err != nil {
			return err
		}

		tr := &entity.Transfer{
			ID:             uuid.NewString(),
			MerchantID:     merchant.ID,
			IdempotencyKey: req.IdempotencyKey,
			PartnerID:      constant.DefaultPartnerID,
			BankCode:       req.BankCode,
			AccountNumber:  req.AccountNumber,
			Amount:         req.Amount,
			Status:         constant.StatusPending,
		}
		if err := u.transferRepo.Create(txCtx, tr); err != nil {
			return err
		}

		mutation := &entity.BalanceMutation{
			MerchantID:    merchant.ID,
			TransferID:    tr.ID,
			Type:          constant.MutationDebit,
			Amount:        tr.Amount,
			BalanceBefore: merchant.Balance,
			BalanceAfter:  balanceAfter,
		}
		if err := u.mutationRepo.Create(txCtx, mutation); err != nil {
			return err
		}

		created = tr
		return nil
	})
	if err != nil {
		if errors.Is(err, transferrepo.ErrDuplicateTransfer) {
			return u.handleDuplicateRace(ctx, idemCacheKey, req.IdempotencyKey)
		}
		return nil, err
	}

	resp := &entity.CreateTransferResponse{
		TransferID: created.ID,
		Status:     created.Status,
		Amount:     created.Amount,
	}

	if err := u.transferRepo.PublishTransfer(ctx, entity.TransferMessage{TransferID: created.ID}); err != nil {
		log.Printf("[api] failed to publish transfer_id=%s to nsq: %v; compensating balance", created.ID, err)
		if compErr := u.refundTransfer(context.Background(), created.ID, "broker_publish_failed"); compErr != nil {
			log.Printf("[api] CRITICAL: failed to compensate transfer_id=%s: %v", created.ID, compErr)
		}
		return nil, fmt.Errorf("%w: %v", ErrBrokerUnavailable, err)
	}

	u.cacheIdempotency(ctx, idemCacheKey, resp)
	return resp, nil
}

func (u *useCase) handleDuplicateRace(ctx context.Context, idemCacheKey, idempotencyKey string) (*entity.CreateTransferResponse, error) {
	existing, err := u.transferRepo.GetByIdempotencyKey(ctx, constant.DefaultMerchantID, idempotencyKey)
	if err != nil {
		return nil, err
	}
	resp := &entity.CreateTransferResponse{
		TransferID: existing.ID,
		Status:     existing.Status,
		Amount:     existing.Amount,
	}
	u.cacheIdempotency(ctx, idemCacheKey, resp)
	return resp, nil
}

// checkIdempotencyPreLock checks Redis first and falls back to a lock-free DB read
// so duplicate/spam retries avoid contending on the merchant row FOR UPDATE lock.
func (u *useCase) checkIdempotencyPreLock(ctx context.Context, idemCacheKey, idempotencyKey string) (*entity.CreateTransferResponse, bool, error) {
	if u.cache != nil {
		var cached entity.CreateTransferResponse
		hit, err := u.cache.GetJSON(ctx, idemCacheKey, &cached)
		if err != nil {
			log.Printf("[api] redis idempotency get failed (falling back to db): %v", err)
		} else if hit {
			return &cached, true, nil
		}
	}

	existing, err := u.transferRepo.GetByIdempotencyKey(ctx, constant.DefaultMerchantID, idempotencyKey)
	if err != nil {
		if errors.Is(err, transferrepo.ErrTransferNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}

	resp := &entity.CreateTransferResponse{
		TransferID: existing.ID,
		Status:     existing.Status,
		Amount:     existing.Amount,
	}
	u.cacheIdempotency(ctx, idemCacheKey, resp)
	return resp, true, nil
}

func (u *useCase) cacheIdempotency(ctx context.Context, key string, resp *entity.CreateTransferResponse) {
	if u.cache == nil {
		return
	}
	if err := u.cache.SetJSON(ctx, key, resp, u.idempotencyTTL); err != nil {
		log.Printf("[api] redis idempotency set failed (ignored): %v", err)
	}
}
