package main

import (
	"encoding/json"
	"github.com/Tango-Rocker/batch-challange/app"
	"github.com/Tango-Rocker/batch-challange/csv"
	"github.com/Tango-Rocker/batch-challange/schema"
	"os"
)

func main() {
	cfg, err := app.LoadConfig()
	if err != nil {
		panic(err)
	}

	p := setupDependencies(err, cfg)

	app.New(cfg, p).Run()
}

func setupDependencies(err error, cfg app.Config) csv.Parser {
	f, err := os.Open(cfg.SchemaPath)
	if err != nil {
		panic(err)
		return nil
	}
	defer f.Close()

	var def schema.CSV
	if err := json.NewDecoder(f).Decode(&def); err != nil {
		panic(err)
	}

	p := csv.NewCSVParser(def)
	return p
}
