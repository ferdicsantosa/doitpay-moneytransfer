package transfernsq

import (
	"context"
	"encoding/json"
	"log"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/entity"
	"github.com/nsqio/go-nsq"
)

func (h *handler) HandleMessage(msg *nsq.Message) error {
	var payload entity.TransferMessage
	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		log.Printf("[nsq] invalid message payload, discarding: %v\n", err)
		return nil // ACK malformed messages so they don't loop infinitely
	}

	shouldRequeue, err := h.usecase.ProcessTransfer(context.Background(), payload.TransferID, msg.Attempts)
	if err != nil && shouldRequeue {
		return err
	}
	return nil
}
