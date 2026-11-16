package worker

import (
	"database/sql"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/config"
)

type Container struct {
	App      *AppModule
	Repos    *RepoModule
	UseCases *UseCaseModule
	Handlers *HandlerModule
}

func BuildContainer(cfg *config.AppConfig, db *sql.DB) *Container {
	app := NewAppModule(cfg, db)
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
