package mutation

import (
	"context"
	"database/sql"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/entity"
)

type Repository interface {
	Create(ctx context.Context, m *entity.BalanceMutation) error
}

type repoModule struct {
	db *sql.DB
}
