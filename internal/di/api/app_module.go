package api

import (
	"database/sql"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/config"
	"github.com/ferdicsantosa/doitpay-moneytransfer/pkg/helper"
	pkgnsq "github.com/ferdicsantosa/doitpay-moneytransfer/pkg/nsq"
	pkgredis "github.com/ferdicsantosa/doitpay-moneytransfer/pkg/redis"
	"github.com/nsqio/go-nsq"
	"github.com/redis/go-redis/v9"
)

type AppModule struct {
	Config    *config.AppConfig
	DB        *sql.DB
	TxManager helper.TxManager
	Publisher pkgnsq.Publisher
	Cache     pkgredis.Cache
}

func NewAppModule(cfg *config.AppConfig, db *sql.DB, producer *nsq.Producer, rdb *redis.Client) *AppModule {
	var cache pkgredis.Cache
	if rdb != nil {
		cache = pkgredis.NewCache(rdb)
	}

	return &AppModule{
		Config:    cfg,
		DB:        db,
		TxManager: helper.NewTxManager(db),
		Publisher: pkgnsq.NewPublisher(producer),
		Cache:     cache,
	}
}
