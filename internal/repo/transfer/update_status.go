package transfer

import (
	"context"
	"errors"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/constant"
)

func (r *repoModule) UpdateStatus(ctx context.Context, transferID string, status constant.TransferStatus, retryCount int, failureReason *string) error {
	// TODO: implement atomic status update
	return errors.New("not implemented")
}
