package transfer

import (
	"context"
	"errors"
)

func (u *useCase) ProcessTransfer(ctx context.Context, transferID string, attempt uint16) (shouldRequeue bool, err error) {
	// TODO: fetch transfer, skip if terminal, mark PROCESSING, call partner, handle success/retry/refund
	return false, errors.New("not implemented")
}
