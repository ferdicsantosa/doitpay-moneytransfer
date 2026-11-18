package merchant

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/entity"
	"github.com/ferdicsantosa/doitpay-moneytransfer/pkg/helper"
)

func (r *repoModule) GetByID(ctx context.Context, merchantID int64) (*entity.Merchant, error) {
	const query = `
		SELECT id, name, balance, created_at, updated_at
		FROM merchants
		WHERE id = $1
	`
	return r.scanMerchant(ctx, query, merchantID)
}

func (r *repoModule) GetByIDForUpdate(ctx context.Context, merchantID int64) (*entity.Merchant, error) {
	const query = `
		SELECT id, name, balance, created_at, updated_at
		FROM merchants
		WHERE id = $1
		FOR UPDATE
	`
	return r.scanMerchant(ctx, query, merchantID)
}

func (r *repoModule) scanMerchant(ctx context.Context, query string, merchantID int64) (*entity.Merchant, error) {
	exec := helper.GetExecutor(ctx, r.db)

	var m entity.Merchant
	err := exec.QueryRowContext(ctx, query, merchantID).Scan(
		&m.ID,
		&m.Name,
		&m.Balance,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMerchantNotFound
		}
		return nil, err
	}
	return &m, nil
}
