package merchant

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/entity"
)

func (r *repoModule) GetByID(ctx context.Context, merchantID int64) (*entity.Merchant, error) {
	const query = `
		SELECT id, name, balance, created_at, updated_at
		FROM merchants
		WHERE id = $1
	`

	var m entity.Merchant
	err := r.db.QueryRowContext(ctx, query, merchantID).Scan(
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
