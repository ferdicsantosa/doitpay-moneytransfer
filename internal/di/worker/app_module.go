package worker

import (
	"database/sql"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/config"
	"github.com/ferdicsantosa/doitpay-moneytransfer/pkg/helper"
	pkgredis "github.com/ferdicsantosa/doitpay-moneytransfer/pkg/redis"
	"github.com/redis/go-redis/v9"
)

type AppModule struct {
	Config    *config.AppConfig
	DB        *sql.DB
	TxManager helper.TxManager
	Locker    pkgredis.Locker
	Cache     pkgredis.Cache
}

func NewAppModule(cfg *config.AppConfig, db *sql.DB, rdb *redis.Client) *AppModule {
	return &AppModule{
		Config:    cfg,
		DB:        db,
		TxManager: helper.NewTxManager(db),
		Locker:    pkgredis.NewLocker(rdb),
		Cache:     pkgredis.NewCache(rdb),
	}
}
