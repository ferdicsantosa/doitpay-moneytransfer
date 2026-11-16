package api

import (
	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/handler/transferhttp"
)

type HandlerModule struct {
	HTTP transferhttp.Handler
}

func NewHandlerModule(usecases *UseCaseModule) *HandlerModule {
	return &HandlerModule{
		HTTP: transferhttp.NewHandler(usecases.Transfer),
	}
}
