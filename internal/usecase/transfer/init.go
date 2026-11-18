package transfer

import (
	"time"

	merchantrepo "github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/merchant"
	mutationrepo "github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/mutation"
	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/partner"
	transferrepo "github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/transfer"
	"github.com/ferdicsantosa/doitpay-moneytransfer/pkg/helper"
	pkgredis "github.com/ferdicsantosa/doitpay-moneytransfer/pkg/redis"
)

func NewAPIUseCase(
	txManager helper.TxManager,
	transferRepo transferrepo.Repository,
	merchantRepo merchantrepo.Repository,
	mutationRepo mutationrepo.Repository,
	cache pkgredis.Cache,
	idempotencyTTL time.Duration,
) APIUseCase {
	return &useCase{
		txManager:      txManager,
		transferRepo:   transferRepo,
		merchantRepo:   merchantRepo,
		mutationRepo:   mutationRepo,
		cache:          cache,
		idempotencyTTL: idempotencyTTL,
	}
}

func NewWorkerUseCase(
	txManager helper.TxManager,
	transferRepo transferrepo.Repository,
	merchantRepo merchantrepo.Repository,
	mutationRepo mutationrepo.Repository,
	partnerRepo partner.Repository,
	locker pkgredis.Locker,
	cache pkgredis.Cache,
	lockTTL time.Duration,
	maxAttempts uint16,
) WorkerUseCase {
	return &useCase{
		txManager:    txManager,
		transferRepo: transferRepo,
		merchantRepo: merchantRepo,
		mutationRepo: mutationRepo,
		partnerRepo:  partnerRepo,
		locker:       locker,
		cache:        cache,
		lockTTL:      lockTTL,
		maxAttempts:  maxAttempts,
	}
}
