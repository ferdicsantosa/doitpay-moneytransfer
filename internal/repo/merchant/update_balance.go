package merchant

import (
	"context"

	"github.com/ferdicsantosa/doitpay-moneytransfer/pkg/helper"
)

func (r *repoModule) UpdateBalance(ctx context.Context, merchantID int64, newBalance int64) error {
	if newBalance < 0 {
		return ErrInsufficientBalance
	}

	const query = `
		UPDATE merchants
		SET balance = $1, updated_at = NOW()
		WHERE id = $2
	`
	exec := helper.GetExecutor(ctx, r.db)
	res, err := exec.ExecContext(ctx, query, newBalance, merchantID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrMerchantNotFound
	}
	return nil
}
