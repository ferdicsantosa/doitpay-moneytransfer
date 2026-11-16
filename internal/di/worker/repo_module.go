package worker

import (
	merchantrepo "github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/merchant"
	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/partner"
	transferrepo "github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/transfer"
)

type RepoModule struct {
	Transfer transferrepo.Repository
	Merchant merchantrepo.Repository
	Partner  partner.Repository
}

func NewRepoModule(app *AppModule) *RepoModule {
	return &RepoModule{
		Transfer: transferrepo.NewRepo(app.DB),
		Merchant: merchantrepo.NewRepo(app.DB),
		Partner:  partner.NewRepo(),
	}
}
