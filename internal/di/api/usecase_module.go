package api

import (
	"time"

	transferuc "github.com/ferdicsantosa/doitpay-moneytransfer/internal/usecase/transfer"
)

type UseCaseModule struct {
	Transfer transferuc.APIUseCase
}

func NewUseCaseModule(app *AppModule, repos *RepoModule) *UseCaseModule {
	idemTTL := time.Duration(app.Config.Redis.IdempotencyTTLSeconds) * time.Second
	if idemTTL <= 0 {
		idemTTL = 5 * time.Minute
	}

	return &UseCaseModule{
		Transfer: transferuc.NewAPIUseCase(
			app.TxManager,
			repos.Transfer,
			repos.Merchant,
			repos.Mutation,
			app.Cache,
			idemTTL,
		),
	}
}
