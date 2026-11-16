package transfer

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/constant"
	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/entity"
)

var (
	ErrTransferNotFound = errors.New("transfer not found")
)

type Repository interface {
	Create(ctx context.Context, transfer *entity.Transfer) error
	GetByID(ctx context.Context, transferID string) (*entity.Transfer, error)
	GetByIDForUpdate(ctx context.Context, transferID string) (*entity.Transfer, error)
	GetByIdempotencyKey(ctx context.Context, merchantID int64, key string) (*entity.Transfer, error)
	UpdateStatus(ctx context.Context, transferID string, status constant.TransferStatus, retryCount int, failureReason *string) error
}

type repoModule struct {
	db *sql.DB
}
