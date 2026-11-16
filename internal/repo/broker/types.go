package broker

import (
	"context"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/entity"
	"github.com/nsqio/go-nsq"
)

type Repository interface {
	PublishTransfer(ctx context.Context, msg entity.TransferMessage) error
}

type repoModule struct {
	producer *nsq.Producer
	topic    string
}
