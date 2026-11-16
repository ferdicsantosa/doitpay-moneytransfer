package transfer

import (
	"context"
	"errors"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/entity"
)

// CreateAndDeductBalance atomically checks idempotency, locks the merchant row,
// deducts the balance, inserts the transfer record, and records a DEBIT mutation.
func (r *repoModule) CreateAndDeductBalance(ctx context.Context, transfer *entity.Transfer) (*entity.Transfer, bool, error) {
	// TODO: implement DB transaction with SELECT ... FOR UPDATE
	return nil, false, errors.New("not implemented")
}
