package api

import (
	"database/sql"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/config"
	"github.com/nsqio/go-nsq"
)

type Container struct {
	App      *AppModule
	Repos    *RepoModule
	UseCases *UseCaseModule
	Handlers *HandlerModule
}

func BuildContainer(cfg *config.AppConfig, db *sql.DB, producer *nsq.Producer) *Container {
	app := NewAppModule(cfg, db, producer)
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
