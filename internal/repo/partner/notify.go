package partner

import (
	"context"
	"log"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/entity"
)

// NotifyDummy simulates sending a notification and merchant webhook callback.
func (r *repoModule) NotifyDummy(ctx context.Context, transfer *entity.Transfer) {
	log.Printf("[dummy-webhook] transfer_id=%s status=%s amount=%d\n", transfer.ID, transfer.Status, transfer.Amount)
}
