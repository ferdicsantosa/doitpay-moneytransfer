package transfer

import (
	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/broker"
	merchantrepo "github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/merchant"
	mutationrepo "github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/mutation"
	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/partner"
	transferrepo "github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/transfer"
	"github.com/ferdicsantosa/doitpay-moneytransfer/pkg/helper"
)

func NewAPIUseCase(
	txManager helper.TxManager,
	transferRepo transferrepo.Repository,
	merchantRepo merchantrepo.Repository,
	mutationRepo mutationrepo.Repository,
	brokerRepo broker.Repository,
) APIUseCase {
	return &useCase{
		txManager:    txManager,
		transferRepo: transferRepo,
		merchantRepo: merchantRepo,
		mutationRepo: mutationRepo,
		brokerRepo:   brokerRepo,
	}
}

func NewWorkerUseCase(
	txManager helper.TxManager,
	transferRepo transferrepo.Repository,
	merchantRepo merchantrepo.Repository,
	mutationRepo mutationrepo.Repository,
	partnerRepo partner.Repository,
	maxAttempts uint16,
) WorkerUseCase {
	return &useCase{
		txManager:    txManager,
		transferRepo: transferRepo,
		merchantRepo: merchantRepo,
		mutationRepo: mutationRepo,
		partnerRepo:  partnerRepo,
		maxAttempts:  maxAttempts,
	}
}
