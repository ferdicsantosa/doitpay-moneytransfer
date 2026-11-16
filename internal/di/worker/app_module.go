package worker

import (
	"database/sql"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/config"
	"github.com/ferdicsantosa/doitpay-moneytransfer/pkg/helper"
)

type AppModule struct {
	Config    *config.AppConfig
	DB        *sql.DB
	TxManager helper.TxManager
}

func NewAppModule(cfg *config.AppConfig, db *sql.DB) *AppModule {
	return &AppModule{
		Config:    cfg,
		DB:        db,
		TxManager: helper.NewTxManager(db),
	}
}
