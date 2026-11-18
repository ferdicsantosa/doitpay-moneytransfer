package worker

import (
	merchantrepo "github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/merchant"
	mutationrepo "github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/mutation"
	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/partner"
	transferrepo "github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/transfer"
)

type RepoModule struct {
	Transfer transferrepo.Repository
	Merchant merchantrepo.Repository
	Mutation mutationrepo.Repository
	Partner  partner.Repository
}

func NewRepoModule(app *AppModule) *RepoModule {
	return &RepoModule{
		Transfer: transferrepo.NewRepo(app.DB, nil, ""),
		Merchant: merchantrepo.NewRepo(app.DB),
		Mutation: mutationrepo.NewRepo(app.DB),
		Partner:  partner.NewRepo(),
	}
}
