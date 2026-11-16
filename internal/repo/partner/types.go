package partner

import (
	"context"
	"errors"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/entity"
)

var (
	ErrRetryable    = errors.New("temporary partner error (retryable)")
	ErrNonRetryable = errors.New("permanent partner error (non-retryable)")
)

type Repository interface {
	Disburse(ctx context.Context, transfer *entity.Transfer) error
	NotifyDummy(ctx context.Context, transfer *entity.Transfer)
}

type repoModule struct{}
