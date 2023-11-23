package app

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/shamank/warehouse-service/internal/handler"
	"github.com/shamank/warehouse-service/internal/repository/postgres"
	"github.com/shamank/warehouse-service/internal/server"
	"github.com/shamank/warehouse-service/internal/service"
	"log/slog"
)

type App struct {
	logger     *slog.Logger
	db         *sql.DB
	httpServer *server.Server
}

func NewApp(logger *slog.Logger, db *sql.DB, httpServer *server.Server) *App {
	return &App{
		logger:     logger,
		db:         db,
		httpServer: httpServer,
	}
}

func (a *App) Run(withTestData bool) error {

	repos := postgres.NewPostgresRepo(a.db, a.logger)

	if withTestData {
		if err := repos.GenerateTestData(); err != nil {
			return err
		}
	}

	services := service.NewService(repos, a.logger)
	handlers := handler.NewHandler(services, a.logger)

	a.httpServer.SetHandler(handlers.InitAPIRoutes())

	return a.httpServer.Start()

}

func (a *App) Stop(ctx context.Context) error {
	return a.httpServer.Stop(ctx)
}
