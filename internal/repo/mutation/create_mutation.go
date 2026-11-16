package mutation

import (
	"context"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/entity"
	"github.com/ferdicsantosa/doitpay-moneytransfer/pkg/helper"
)

func (r *repoModule) Create(ctx context.Context, m *entity.BalanceMutation) error {
	const query = `
		INSERT INTO balance_mutations (
			merchant_id, transfer_id, type, amount, balance_before, balance_after
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	exec := helper.GetExecutor(ctx, r.db)
	return exec.QueryRowContext(
		ctx,
		query,
		m.MerchantID,
		m.TransferID,
		m.Type,
		m.Amount,
		m.BalanceBefore,
		m.BalanceAfter,
	).Scan(&m.ID, &m.CreatedAt)
}
