package api

import (
	"database/sql"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/config"
	"github.com/nsqio/go-nsq"
)

type AppModule struct {
	Config   *config.AppConfig
	DB       *sql.DB
	Producer *nsq.Producer
}

func NewAppModule(cfg *config.AppConfig, db *sql.DB, producer *nsq.Producer) *AppModule {
	return &AppModule{
		Config:   cfg,
		DB:       db,
		Producer: producer,
	}
}
