package main

import (
	"context"
	"github.com/Tango-Rocker/batch-challange/app"
	"github.com/Tango-Rocker/batch-challange/business"
	"github.com/Tango-Rocker/batch-challange/csv"
	"github.com/Tango-Rocker/batch-challange/db"
	"github.com/Tango-Rocker/batch-challange/utils"
	"log/slog"
)

//TODO: make execution-id unique and add a model for it

func main() {

	ctx := context.Background()

	logger := slog.Default().With(
		slog.String("app", "batch-challenge"),
		slog.String("version", "0.0.1"),
	)

	/*mailCfg, err := utils.LoadEnvConfig[business.MailConfig]()
	if err != nil {
		panic(err)
	}
	mailService := business.NewEmailService(mailCfg, logger)
	*/
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
	writer := setupWriter(repository, logger)

	app.New(appCfg, parser, writer, logger).Run(ctx)
}

func setupRepository(dbCfg *db.Config, logger *slog.Logger) *db.Repository {
	repository, err := db.NewRepository(dbCfg, logger)
	if err != nil {
		panic(err)
	}
	return repository
}

func setupWriter(repository *db.Repository, logger *slog.Logger) *business.Writer {
	bsCfg, err := utils.LoadEnvConfig[business.WriterConfig]()
	if err != nil {
		panic(err)
	}

	worker := business.NewWriter(bsCfg, repository, logger)
	return worker
}
