package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/shamank/warehouse-service/internal/app"
	"github.com/shamank/warehouse-service/internal/config"
	"github.com/shamank/warehouse-service/internal/server"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	var configPath string    // путь к файлу конфигурации
	var migrationPath string // путь к папке с миграциями

	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.StringVar(&migrationPath, "migrate", "", "migrate database")
	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH") // если не задан путь к конфигурации, то берем из переменной окружения
	}

	if migrationPath == "" {
		migrationPath = os.Getenv("MIGRATION_PATH") // если не задан путь к миграциям, то берем из переменной окружения
	}

	cfg := config.InitConfig(configPath)

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Database, cfg.Postgres.SSLMode),
	)
	if err != nil {
		return
	}
	defer db.Close()

	logger := initLogger("debug")

	if migrationPath != "" {
		if err := checkMigrations(db, migrationPath); err != nil {
			logger.Error(err.Error())
			return
		}
	}

	serv := server.NewServer(cfg.HTTP)

	application := app.NewApp(logger, db, serv)
	go func() {
		logger.Info("warehouse service started!")
		if err := application.Run(cfg.InsertTestData); err != nil {
			logger.Error(err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second
	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := application.Stop(ctx); err != nil {
		return
	}

	logger.Info("warehouse service stopped")
}

func checkMigrations(db *sql.DB, migrationPath string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationPath,
		"postgres",
		driver,
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(err)
	}
	return nil
}

func initLogger(levelString string) *slog.Logger {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	return logger
}
