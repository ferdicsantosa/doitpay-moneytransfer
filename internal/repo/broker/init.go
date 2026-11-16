package broker

import "github.com/nsqio/go-nsq"

func NewRepo(producer *nsq.Producer, topic string) Repository {
	return &repoModule{
		producer: producer,
		topic:    topic,
	}
}
