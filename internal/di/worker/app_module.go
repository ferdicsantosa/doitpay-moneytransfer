package worker

import (
	"database/sql"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/config"
)

type AppModule struct {
	Config *config.AppConfig
	DB     *sql.DB
}

func NewAppModule(cfg *config.AppConfig, db *sql.DB) *AppModule {
	return &AppModule{
		Config: cfg,
		DB:     db,
	}
}
