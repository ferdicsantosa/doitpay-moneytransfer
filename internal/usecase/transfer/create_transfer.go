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
	var (
		created     *entity.Transfer
		isDuplicate bool
	)

	err := u.txManager.WithTx(ctx, func(txCtx context.Context) error {
		// Lock merchant row first to serialize concurrent deductions and idempotency checks.
		merchant, err := u.merchantRepo.GetByIDForUpdate(txCtx, constant.DefaultMerchantID)
		if err != nil {
			return err
		}

		existing, err := u.transferRepo.GetByIdempotencyKey(txCtx, merchant.ID, req.IdempotencyKey)
		if err != nil && !errors.Is(err, transferrepo.ErrTransferNotFound) {
			return err
		}
		if existing != nil {
			created = existing
			isDuplicate = true
			return nil
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
		return nil, err
	}

	if isDuplicate {
		return &entity.CreateTransferResponse{
			TransferID: created.ID,
			Status:     created.Status,
			Amount:     created.Amount,
		}, nil
	}

	if err := u.brokerRepo.PublishTransfer(ctx, entity.TransferMessage{TransferID: created.ID}); err != nil {
		log.Printf("[api] failed to publish transfer_id=%s to nsq: %v; compensating balance", created.ID, err)
		if compErr := u.refundTransfer(context.Background(), created.ID, "broker_publish_failed"); compErr != nil {
			log.Printf("[api] CRITICAL: failed to compensate transfer_id=%s: %v", created.ID, compErr)
		}
		return nil, fmt.Errorf("%w: %v", ErrBrokerUnavailable, err)
	}

	return &entity.CreateTransferResponse{
		TransferID: created.ID,
		Status:     created.Status,
		Amount:     created.Amount,
	}, nil
}
