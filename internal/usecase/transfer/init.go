package transfer

import (
	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/broker"
	merchantrepo "github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/merchant"
	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/partner"
	transferrepo "github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/transfer"
)

func NewAPIUseCase(
	transferRepo transferrepo.Repository,
	merchantRepo merchantrepo.Repository,
	brokerRepo broker.Repository,
) APIUseCase {
	return &useCase{
		transferRepo: transferRepo,
		merchantRepo: merchantRepo,
		brokerRepo:   brokerRepo,
	}
}

func NewWorkerUseCase(
	transferRepo transferrepo.Repository,
	merchantRepo merchantrepo.Repository,
	partnerRepo partner.Repository,
	maxAttempts uint16,
) WorkerUseCase {
	return &useCase{
		transferRepo: transferRepo,
		merchantRepo: merchantRepo,
		partnerRepo:  partnerRepo,
		maxAttempts:  maxAttempts,
	}
}
