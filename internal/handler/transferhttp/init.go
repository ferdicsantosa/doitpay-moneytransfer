package transferhttp

import transferuc "github.com/ferdicsantosa/doitpay-moneytransfer/internal/usecase/transfer"

func NewHandler(uc transferuc.APIUseCase) Handler {
	return &handler{
		usecase: uc,
	}
}
