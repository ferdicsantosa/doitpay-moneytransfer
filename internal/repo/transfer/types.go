package transfer

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/constant"
	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/entity"
)

var (
	ErrInsufficientBalance = errors.New("insufficient merchant balance")
	ErrMerchantNotFound    = errors.New("merchant not found")
	ErrTransferNotFound    = errors.New("transfer not found")
)

type Repository interface {
	CreateAndDeductBalance(ctx context.Context, transfer *entity.Transfer) (*entity.Transfer, bool, error)
	GetByID(ctx context.Context, transferID string) (*entity.Transfer, error)
	UpdateStatus(ctx context.Context, transferID string, status constant.TransferStatus, retryCount int, failureReason *string) error
	MarkFailedAndRefund(ctx context.Context, transferID string, reason string) error
}

type repoModule struct {
	db *sql.DB
}
