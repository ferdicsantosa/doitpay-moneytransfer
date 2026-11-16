package api

import (
	transferuc "github.com/ferdicsantosa/doitpay-moneytransfer/internal/usecase/transfer"
)

type UseCaseModule struct {
	Transfer transferuc.APIUseCase
}

func NewUseCaseModule(app *AppModule, repos *RepoModule) *UseCaseModule {
	return &UseCaseModule{
		Transfer: transferuc.NewAPIUseCase(
			app.TxManager,
			repos.Transfer,
			repos.Merchant,
			repos.Mutation,
			repos.Broker,
		),
	}
}
