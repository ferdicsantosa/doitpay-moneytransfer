package transferhttp

import (
	"net/http"

	transferuc "github.com/ferdicsantosa/doitpay-moneytransfer/internal/usecase/transfer"
)

type Handler interface {
	CreateTransfer(w http.ResponseWriter, r *http.Request)
	Health(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	usecase transferuc.APIUseCase
}
