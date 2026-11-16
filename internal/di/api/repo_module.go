package api

import (
	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/broker"
	merchantrepo "github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/merchant"
	mutationrepo "github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/mutation"
	transferrepo "github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/transfer"
)

type RepoModule struct {
	Transfer transferrepo.Repository
	Merchant merchantrepo.Repository
	Mutation mutationrepo.Repository
	Broker   broker.Repository
}

func NewRepoModule(app *AppModule) *RepoModule {
	return &RepoModule{
		Transfer: transferrepo.NewRepo(app.DB),
		Merchant: merchantrepo.NewRepo(app.DB),
		Mutation: mutationrepo.NewRepo(app.DB),
		Broker:   broker.NewRepo(app.Producer, app.Config.NSQ.Topic),
	}
}
