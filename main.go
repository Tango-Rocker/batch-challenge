package main

import (
	"context"
	"github.com/Tango-Rocker/batch-challange/app"
	"github.com/Tango-Rocker/batch-challange/csv"
	"github.com/Tango-Rocker/batch-challange/db"
	"log/slog"
)

func main() {
	//this should be parser config
	ctx := context.Background()

	appCfg, err := app.LoadConfig()
	if err != nil {
		panic(err)
	}

	logger := slog.Default().With(
		slog.String("app", "batch-challenge"),
		slog.String("version", "0.0.1"),
	)

	parser := csv.NewCSVParser(appCfg.SchemaPath, logger)

	dbCfg, err := db.LoadConfig()
	if err != nil {
		panic(err)
	}

	dbService, err := db.NewService(dbCfg, ctx, logger)
	if err != nil {
		panic(err)
	}

	app.New(appCfg, parser, dbService).Run(ctx)
}
