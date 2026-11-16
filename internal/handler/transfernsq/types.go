package transfernsq

import (
	transferuc "github.com/ferdicsantosa/doitpay-moneytransfer/internal/usecase/transfer"
	"github.com/nsqio/go-nsq"
)

type Handler interface {
	nsq.Handler
}

type handler struct {
	usecase transferuc.WorkerUseCase
}
