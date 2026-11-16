package merchant

import (
	"context"
	"database/sql"
	"errors"
)

// UpdateBalance atomically adjusts a merchant's balance by delta (negative to deduct, positive to credit).
func (r *repoModule) UpdateBalance(ctx context.Context, merchantID int64, delta int64) (int64, int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, 0, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var balanceBefore int64
	const lockQuery = `SELECT balance FROM merchants WHERE id = $1 FOR UPDATE`
	if err := tx.QueryRowContext(ctx, lockQuery, merchantID).Scan(&balanceBefore); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, 0, ErrMerchantNotFound
		}
		return 0, 0, err
	}

	balanceAfter := balanceBefore + delta
	if balanceAfter < 0 {
		return 0, 0, ErrInsufficientBalance
	}

	const updateQuery = `
		UPDATE merchants
		SET balance = $1, updated_at = NOW()
		WHERE id = $2
	`
	if _, err := tx.ExecContext(ctx, updateQuery, balanceAfter, merchantID); err != nil {
		return 0, 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, 0, err
	}

	return balanceBefore, balanceAfter, nil
}
