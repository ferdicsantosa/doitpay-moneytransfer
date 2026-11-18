package transfer

import (
	"context"
	"errors"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/entity"
)

func (r *repoModule) PublishTransfer(ctx context.Context, msg entity.TransferMessage) error {
	if r.publisher == nil {
		return errors.New("nsq publisher is not configured")
	}
	return r.publisher.Publish(ctx, r.topic, msg)
}
