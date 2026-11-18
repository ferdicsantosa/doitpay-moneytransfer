package worker

import (
	"time"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/handler/transfernsq"
)

type HandlerModule struct {
	NSQ transfernsq.Handler
}

func NewHandlerModule(app *AppModule, usecases *UseCaseModule) *HandlerModule {
	baseDelay := time.Duration(app.Config.NSQ.RequeueDelaySeconds) * time.Second
	maxDelay := time.Duration(app.Config.NSQ.MaxRequeueDelaySeconds) * time.Second
	return &HandlerModule{
		NSQ: transfernsq.NewHandler(usecases.Transfer, baseDelay, maxDelay),
	}
}
