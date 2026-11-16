package transfer

import (
	"context"
	"errors"
)

// MarkFailedAndRefund atomically marks a non-terminal transfer as FAILED,
// refunds the merchant balance, and records a CREDIT mutation.
func (r *repoModule) MarkFailedAndRefund(ctx context.Context, transferID string, reason string) error {
	// TODO: implement idempotent refund inside DB transaction
	return errors.New("not implemented")
}
