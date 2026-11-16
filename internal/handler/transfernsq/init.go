package transfernsq

import transferuc "github.com/ferdicsantosa/doitpay-moneytransfer/internal/usecase/transfer"

func NewHandler(uc transferuc.WorkerUseCase) Handler {
	return &handler{
		usecase: uc,
	}
}
