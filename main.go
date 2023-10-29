package main

import (
	"context"
	"github.com/Tango-Rocker/batch-challange/app"
	"github.com/Tango-Rocker/batch-challange/bussines"
	"github.com/Tango-Rocker/batch-challange/csv"
	"github.com/Tango-Rocker/batch-challange/db"
	"github.com/Tango-Rocker/batch-challange/utils"
	"log/slog"
)

func main() {
	ctx := context.Background()

	logger := slog.Default().With(
		slog.String("app", "batch-challenge"),
		slog.String("version", "0.0.1"),
	)

	appCfg, err := utils.LoadEnvConfig[app.Config]()
	if err != nil {
		panic(err)
	}

	dbCfg, err := utils.LoadEnvConfig[db.Config]()
	if err != nil {
		panic(err)
	}

	parser := csv.NewCSVParser(appCfg.SchemaPath, logger)
	repository := setupRepository(dbCfg, logger)
	worker := setupInsertWorker(repository, logger)

	app.New(appCfg, parser, worker).Run(ctx)
}

func setupRepository(dbCfg *db.Config, logger *slog.Logger) *db.Repository {
	repository, err := db.NewRepository(dbCfg, logger)
	if err != nil {
		panic(err)
	}
	return repository
}

func setupInsertWorker(repository *db.Repository, logger *slog.Logger) *bussines.Worker {
	bsCfg, err := utils.LoadEnvConfig[bussines.Config]()
	if err != nil {
		panic(err)
	}

	worker := bussines.NewWorker(bsCfg, repository, logger)
	return worker
}
