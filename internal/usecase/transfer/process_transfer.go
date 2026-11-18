package transfer

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/constant"
	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/partner"
	transferrepo "github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/transfer"
)

func (u *useCase) ProcessTransfer(ctx context.Context, transferID string, attempt uint16) (bool, error) {
	lockKey := "lock:transfer:" + transferID
	token, acquired, err := u.locker.Acquire(ctx, lockKey, u.lockTTL)
	if err != nil {
		log.Printf("[worker] transfer_id=%s failed to acquire lock: %v", transferID, err)
		return true, err
	}
	if !acquired {
		log.Printf("[worker] transfer_id=%s already being processed by another worker, skipping", transferID)
		return false, nil
	}
	defer func() {
		if relErr := u.locker.Release(context.Background(), lockKey, token); relErr != nil {
			log.Printf("[worker] transfer_id=%s failed to release lock: %v", transferID, relErr)
		}
	}()

	tr, err := u.transferRepo.GetByID(ctx, transferID)
	if err != nil {
		if errors.Is(err, transferrepo.ErrTransferNotFound) {
			log.Printf("[worker] transfer_id=%s not found, discarding message", transferID)
			return false, nil
		}
		return true, err
	}

	if tr.Status.IsTerminal() {
		log.Printf("[worker] transfer_id=%s already in terminal status=%s, skipping", tr.ID, tr.Status)
		return false, nil
	}

	retryCount := 0
	if attempt > 1 {
		retryCount = int(attempt - 1)
	}

	disburseErr := u.partnerRepo.Disburse(ctx, tr)
	if disburseErr == nil {
		if err := u.transferRepo.UpdateStatus(ctx, tr.ID, constant.StatusSuccess, retryCount, nil); err != nil {
			return true, err
		}
		u.deleteIdempotencyCache(ctx, tr.MerchantID, tr.IdempotencyKey)
		tr.Status = constant.StatusSuccess
		u.partnerRepo.NotifyDummy(ctx, tr)
		return false, nil
	}

	if errors.Is(disburseErr, partner.ErrRetryable) && attempt < u.maxAttempts {
		log.Printf("[worker] transfer_id=%s attempt=%d/%d retryable error, requeuing: %v", tr.ID, attempt, u.maxAttempts, disburseErr)
		return true, disburseErr
	}

	log.Printf("[worker] transfer_id=%s attempt=%d/%d permanent failure, refunding: %v", tr.ID, attempt, u.maxAttempts, disburseErr)
	if err := u.refundTransfer(ctx, tr.ID, disburseErr.Error()); err != nil {
		return true, err
	}

	u.deleteIdempotencyCache(ctx, tr.MerchantID, tr.IdempotencyKey)
	tr.Status = constant.StatusFailed
	u.partnerRepo.NotifyDummy(ctx, tr)
	return false, nil
}

func (u *useCase) deleteIdempotencyCache(ctx context.Context, merchantID int64, idempotencyKey string) {
	if u.cache == nil {
		return
	}
	key := fmt.Sprintf("idem:transfer:%d:%s", merchantID, idempotencyKey)
	if err := u.cache.Delete(ctx, key); err != nil {
		log.Printf("[worker] failed to delete idempotency cache key=%s (ignored): %v", key, err)
	}
}
