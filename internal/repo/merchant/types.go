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
	UpdateBalance(ctx context.Context, merchantID int64, delta int64) (balanceBefore int64, balanceAfter int64, err error)
}

type repoModule struct {
	db *sql.DB
}
