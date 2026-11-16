package api

import (
	transferuc "github.com/ferdicsantosa/doitpay-moneytransfer/internal/usecase/transfer"
)

type UseCaseModule struct {
	Transfer transferuc.APIUseCase
}

func NewUseCaseModule(repos *RepoModule) *UseCaseModule {
	return &UseCaseModule{
		Transfer: transferuc.NewAPIUseCase(
			repos.Transfer,
			repos.Merchant,
			repos.Broker,
		),
	}
}
