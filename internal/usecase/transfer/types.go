package transfer

import (
	"context"
	"errors"
	"time"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/entity"
	merchantrepo "github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/merchant"
	mutationrepo "github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/mutation"
	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/partner"
	transferrepo "github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/transfer"
	"github.com/ferdicsantosa/doitpay-moneytransfer/pkg/helper"
	pkgredis "github.com/ferdicsantosa/doitpay-moneytransfer/pkg/redis"
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
	txManager      helper.TxManager
	transferRepo   transferrepo.Repository
	merchantRepo   merchantrepo.Repository
	mutationRepo   mutationrepo.Repository
	partnerRepo    partner.Repository
	locker         pkgredis.Locker
	cache          pkgredis.Cache
	lockTTL        time.Duration
	idempotencyTTL time.Duration
	maxAttempts    uint16
}
