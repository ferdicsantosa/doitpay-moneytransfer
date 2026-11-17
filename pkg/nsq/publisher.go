package nsq

import (
	"context"
	"encoding/json"

	"github.com/nsqio/go-nsq"
)

type Publisher interface {
	Publish(ctx context.Context, topic string, payload any) error
}

type publisher struct {
	producer *nsq.Producer
}

func NewPublisher(producer *nsq.Producer) Publisher {
	return &publisher{producer: producer}
}

func (p *publisher) Publish(ctx context.Context, topic string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return p.producer.Publish(topic, data)
}
