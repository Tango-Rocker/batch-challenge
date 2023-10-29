package main

import (
	"github.com/Tango-Rocker/batch-challange/app"
	"github.com/Tango-Rocker/batch-challange/csv"
	"github.com/Tango-Rocker/batch-challange/db"
)

func main() {
	//this should be parser config
	appCfg, err := app.LoadConfig()
	if err != nil {
		panic(err)
	}
	parser := csv.NewCSVParser(appCfg.SchemaPath)

	dbCfg, err := db.LoadConfig()
	if err != nil {
		panic(err)
	}

	dbService, err := db.NewService(dbCfg)
	if err != nil {
		panic(err)
	}

	app.New(appCfg, parser, dbService).Run()
}
