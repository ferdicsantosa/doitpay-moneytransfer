package transfer

import (
	"context"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/entity"
	"github.com/ferdicsantosa/doitpay-moneytransfer/pkg/helper"
)

func (r *repoModule) Create(ctx context.Context, transfer *entity.Transfer) error {
	const query = `
		INSERT INTO transfers (
			id, merchant_id, idempotency_key, partner_id,
			bank_code, account_number, amount, status, retry_count
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at, updated_at
	`

	exec := helper.GetExecutor(ctx, r.db)
	return exec.QueryRowContext(
		ctx,
		query,
		transfer.ID,
		transfer.MerchantID,
		transfer.IdempotencyKey,
		transfer.PartnerID,
		transfer.BankCode,
		transfer.AccountNumber,
		transfer.Amount,
		transfer.Status,
		transfer.RetryCount,
	).Scan(&transfer.CreatedAt, &transfer.UpdatedAt)
}
