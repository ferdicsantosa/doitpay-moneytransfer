package worker

import (
	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/handler/transfernsq"
)

type HandlerModule struct {
	NSQ transfernsq.Handler
}

func NewHandlerModule(usecases *UseCaseModule) *HandlerModule {
	return &HandlerModule{
		NSQ: transfernsq.NewHandler(usecases.Transfer),
	}
}
