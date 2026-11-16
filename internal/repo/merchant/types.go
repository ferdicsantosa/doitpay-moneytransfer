package merchant

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/entity"
)

var (
	ErrMerchantNotFound    = errors.New("merchant not found")
	ErrInsufficientBalance = errors.New("insufficient merchant balance")
)

type Repository interface {
	GetByID(ctx context.Context, merchantID int64) (*entity.Merchant, error)
	GetByIDForUpdate(ctx context.Context, merchantID int64) (*entity.Merchant, error)
	UpdateBalance(ctx context.Context, merchantID int64, newBalance int64) error
}

type repoModule struct {
	db *sql.DB
}
