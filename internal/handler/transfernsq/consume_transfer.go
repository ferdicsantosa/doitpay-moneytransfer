package transfernsq

import (
	"context"
	"encoding/json"
	"log"
	"math/rand/v2"
	"time"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/entity"
	"github.com/nsqio/go-nsq"
)

func (h *handler) HandleMessage(msg *nsq.Message) error {
	var payload entity.TransferMessage
	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		log.Printf("[nsq] invalid message payload, discarding: %v\n", err)
		return nil // ACK malformed messages so they don't loop infinitely
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	shouldRequeue, err := h.usecase.ProcessTransfer(ctx, payload.TransferID, msg.Attempts)
	if err != nil && shouldRequeue {
		delay := h.backoffWithJitter(msg.Attempts)
		log.Printf("[nsq] requeuing transfer_id=%s attempt=%d delay=%s", payload.TransferID, msg.Attempts, delay)
		msg.RequeueWithoutBackoff(delay)
		return nil
	}
	return nil
}

func (h *handler) backoffWithJitter(attempt uint16) time.Duration {
	shift := 0
	if attempt > 1 {
		shift = int(attempt - 1)
	}
	if shift > 10 {
		shift = 10
	}

	delay := h.baseRequeueDelay * time.Duration(1<<shift)
	if delay > h.maxRequeueDelay {
		delay = h.maxRequeueDelay
	}

	// Add up to 50% positive jitter to prevent thundering herd retries.
	jitter := time.Duration(rand.Int64N(int64(delay/2) + 1))
	res := delay + jitter
	if res > h.maxRequeueDelay {
		return h.maxRequeueDelay
	}
	return res
}
