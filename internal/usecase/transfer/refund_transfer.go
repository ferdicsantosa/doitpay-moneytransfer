package transfer

import (
	"context"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/constant"
	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/entity"
)

// refundTransfer atomically marks a non-terminal transfer as FAILED and credits the merchant balance.
func (u *useCase) refundTransfer(ctx context.Context, transferID string, reason string) error {
	return u.txManager.WithTx(ctx, func(txCtx context.Context) error {
		tr, err := u.transferRepo.GetByIDForUpdate(txCtx, transferID)
		if err != nil {
			return err
		}
		if tr.Status.IsTerminal() {
			return nil
		}

		if err := u.transferRepo.UpdateStatus(txCtx, tr.ID, constant.StatusFailed, tr.RetryCount, &reason); err != nil {
			return err
		}

		merchant, err := u.merchantRepo.GetByIDForUpdate(txCtx, tr.MerchantID)
		if err != nil {
			return err
		}

		balanceAfter := merchant.Balance + tr.Amount
		if err := u.merchantRepo.UpdateBalance(txCtx, merchant.ID, balanceAfter); err != nil {
			return err
		}

		return u.mutationRepo.Create(txCtx, &entity.BalanceMutation{
			MerchantID:    merchant.ID,
			TransferID:    tr.ID,
			Type:          constant.MutationCredit,
			Amount:        tr.Amount,
			BalanceBefore: merchant.Balance,
			BalanceAfter:  balanceAfter,
		})
	})
}
