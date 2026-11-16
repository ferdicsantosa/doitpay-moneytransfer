package transfer

import (
	"context"
	"errors"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/entity"
)

func (u *useCase) CreateTransfer(ctx context.Context, req entity.CreateTransferRequest) (*entity.CreateTransferResponse, error) {
	// TODO: validate request, deduct balance + create transfer in DB, publish to NSQ (compensate on publish error)
	return nil, errors.New("not implemented")
}
