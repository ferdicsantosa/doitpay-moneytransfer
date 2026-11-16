package transfer

import (
	"context"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/entity"
	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/broker"
	merchantrepo "github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/merchant"
	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/partner"
	transferrepo "github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/transfer"
)

type APIUseCase interface {
	CreateTransfer(ctx context.Context, req entity.CreateTransferRequest) (*entity.CreateTransferResponse, error)
}

type WorkerUseCase interface {
	ProcessTransfer(ctx context.Context, transferID string, attempt uint16) (shouldRequeue bool, err error)
}

type useCase struct {
	transferRepo transferrepo.Repository
	merchantRepo merchantrepo.Repository
	brokerRepo   broker.Repository
	partnerRepo  partner.Repository
	maxAttempts  uint16
}
