package worker

import (
	transferuc "github.com/ferdicsantosa/doitpay-moneytransfer/internal/usecase/transfer"
)

type UseCaseModule struct {
	Transfer transferuc.WorkerUseCase
}

func NewUseCaseModule(app *AppModule, repos *RepoModule) *UseCaseModule {
	return &UseCaseModule{
		Transfer: transferuc.NewWorkerUseCase(
			app.TxManager,
			repos.Transfer,
			repos.Merchant,
			repos.Mutation,
			repos.Partner,
			app.Config.NSQ.MaxAttempts,
		),
	}
}
