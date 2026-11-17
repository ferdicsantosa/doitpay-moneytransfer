package api

import (
	merchantrepo "github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/merchant"
	mutationrepo "github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/mutation"
	transferrepo "github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/transfer"
)

type RepoModule struct {
	Transfer transferrepo.Repository
	Merchant merchantrepo.Repository
	Mutation mutationrepo.Repository
}

func NewRepoModule(app *AppModule) *RepoModule {
	return &RepoModule{
		Transfer: transferrepo.NewRepo(app.DB, app.Publisher, app.Config.NSQ.Topic),
		Merchant: merchantrepo.NewRepo(app.DB),
		Mutation: mutationrepo.NewRepo(app.DB),
	}
}
