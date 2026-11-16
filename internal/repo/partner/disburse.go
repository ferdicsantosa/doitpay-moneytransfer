package partner

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"time"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/entity"
)

// Disburse simulates calling an external mock payment provider with random delay and outcome.
func (r *repoModule) Disburse(ctx context.Context, transfer *entity.Transfer) error {
	// Simulate network/provider latency between 100ms and 800ms.
	delay := time.Duration(100+rand.IntN(700)) * time.Millisecond

	select {
	case <-time.After(delay):
	case <-ctx.Done():
		return ctx.Err()
	}

	// Simulate outcome: 70% success, 20% retryable error, 10% non-retryable error.
	roll := rand.IntN(100)
	switch {
	case roll < 70:
		log.Printf("[mock-partner] transfer_id=%s status=SUCCESS delay=%s", transfer.ID, delay)
		return nil
	case roll < 90:
		log.Printf("[mock-partner] transfer_id=%s status=RETRYABLE_ERROR delay=%s", transfer.ID, delay)
		return fmt.Errorf("%w: partner gateway timeout", ErrRetryable)
	default:
		log.Printf("[mock-partner] transfer_id=%s status=NON_RETRYABLE_ERROR delay=%s", transfer.ID, delay)
		return fmt.Errorf("%w: invalid destination bank account", ErrNonRetryable)
	}
}
