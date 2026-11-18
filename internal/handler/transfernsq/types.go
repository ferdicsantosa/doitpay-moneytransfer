package transfernsq

import (
	"context"
	"time"

	"github.com/nsqio/go-nsq"
)

type Handler interface {
	nsq.Handler
}

type TransferUsecase interface {
	ProcessTransfer(ctx context.Context, transferID string, attempt uint16) (shouldRequeue bool, err error)
}

type handler struct {
	usecase          TransferUsecase
	baseRequeueDelay time.Duration
	maxRequeueDelay  time.Duration
}
