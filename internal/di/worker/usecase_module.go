package worker

import (
	"time"

	transferuc "github.com/ferdicsantosa/doitpay-moneytransfer/internal/usecase/transfer"
)

type UseCaseModule struct {
	Transfer transferuc.WorkerUseCase
}

func NewUseCaseModule(app *AppModule, repos *RepoModule) *UseCaseModule {
	lockTTL := time.Duration(app.Config.Redis.LockTTLSeconds) * time.Second
	if lockTTL <= 0 {
		lockTTL = 30 * time.Second
	}

	return &UseCaseModule{
		Transfer: transferuc.NewWorkerUseCase(
			app.TxManager,
			repos.Transfer,
			repos.Merchant,
			repos.Mutation,
			repos.Partner,
			app.Locker,
			app.Cache,
			lockTTL,
			app.Config.NSQ.MaxAttempts,
		),
	}
}
