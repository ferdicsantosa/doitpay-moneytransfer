package worker

import (
	"database/sql"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/config"
	"github.com/redis/go-redis/v9"
)

type Container struct {
	App      *AppModule
	Repos    *RepoModule
	UseCases *UseCaseModule
	Handlers *HandlerModule
}

func BuildContainer(cfg *config.AppConfig, db *sql.DB, rdb *redis.Client) *Container {
	app := NewAppModule(cfg, db, rdb)
	repos := NewRepoModule(app)
	usecases := NewUseCaseModule(app, repos)
	handlers := NewHandlerModule(usecases)

	return &Container{
		App:      app,
		Repos:    repos,
		UseCases: usecases,
		Handlers: handlers,
	}
}
