package broker

import (
	"context"
	"encoding/json"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/entity"
)

func (r *repoModule) PublishTransfer(ctx context.Context, msg entity.TransferMessage) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return r.producer.Publish(r.topic, payload)
}
