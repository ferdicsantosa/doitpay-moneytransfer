package transfer

import (
	"context"
	"errors"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/entity"
	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/broker"
	merchantrepo "github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/merchant"
	mutationrepo "github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/mutation"
	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/partner"
	transferrepo "github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/transfer"
	"github.com/ferdicsantosa/doitpay-moneytransfer/pkg/helper"
)

var (
	ErrBrokerUnavailable = errors.New("message broker unavailable")
)

type APIUseCase interface {
	CreateTransfer(ctx context.Context, req entity.CreateTransferRequest) (*entity.CreateTransferResponse, error)
}

type WorkerUseCase interface {
	ProcessTransfer(ctx context.Context, transferID string, attempt uint16) (shouldRequeue bool, err error)
}

type useCase struct {
	txManager    helper.TxManager
	transferRepo transferrepo.Repository
	merchantRepo merchantrepo.Repository
	mutationRepo mutationrepo.Repository
	brokerRepo   broker.Repository
	partnerRepo  partner.Repository
	maxAttempts  uint16
}
