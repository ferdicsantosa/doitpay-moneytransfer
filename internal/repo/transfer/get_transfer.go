package transfer

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/entity"
	"github.com/ferdicsantosa/doitpay-moneytransfer/pkg/helper"
)

func (r *repoModule) GetByID(ctx context.Context, transferID string) (*entity.Transfer, error) {
	const query = `
		SELECT id, merchant_id, idempotency_key, partner_id,
		       bank_code, account_number, amount, status,
		       retry_count, failure_reason, created_at, updated_at
		FROM transfers
		WHERE id = $1
	`
	return r.queryOne(ctx, query, transferID)
}

func (r *repoModule) GetByIDForUpdate(ctx context.Context, transferID string) (*entity.Transfer, error) {
	const query = `
		SELECT id, merchant_id, idempotency_key, partner_id,
		       bank_code, account_number, amount, status,
		       retry_count, failure_reason, created_at, updated_at
		FROM transfers
		WHERE id = $1
		FOR UPDATE
	`
	return r.queryOne(ctx, query, transferID)
}

func (r *repoModule) GetByIdempotencyKey(ctx context.Context, merchantID int64, key string) (*entity.Transfer, error) {
	const query = `
		SELECT id, merchant_id, idempotency_key, partner_id,
		       bank_code, account_number, amount, status,
		       retry_count, failure_reason, created_at, updated_at
		FROM transfers
		WHERE merchant_id = $1 AND idempotency_key = $2
	`
	return r.queryOne(ctx, query, merchantID, key)
}

func (r *repoModule) queryOne(ctx context.Context, query string, args ...any) (*entity.Transfer, error) {
	exec := helper.GetExecutor(ctx, r.db)
	t, err := scanTransfer(exec.QueryRowContext(ctx, query, args...))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTransferNotFound
		}
		return nil, err
	}
	return t, nil
}

func scanTransfer(row *sql.Row) (*entity.Transfer, error) {
	var t entity.Transfer
	var failureReason sql.NullString

	err := row.Scan(
		&t.ID,
		&t.MerchantID,
		&t.IdempotencyKey,
		&t.PartnerID,
		&t.BankCode,
		&t.AccountNumber,
		&t.Amount,
		&t.Status,
		&t.RetryCount,
		&failureReason,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if failureReason.Valid {
		t.FailureReason = &failureReason.String
	}
	return &t, nil
}
