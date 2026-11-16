package transfer

import (
	"context"
	"errors"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/entity"
)

func (r *repoModule) GetByID(ctx context.Context, transferID string) (*entity.Transfer, error) {
	// TODO: implement query by transfer ID
	return nil, errors.New("not implemented")
}
