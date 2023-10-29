package app

import (
	"context"
	"fmt"
	"github.com/Tango-Rocker/batch-challange/csv"
	"github.com/Tango-Rocker/batch-challange/db"
	"os"
)

// TODO: i should refactor this a bit more with reader/writer interfaces
// dont really care what kind of io is it really

type Application struct {
	Config
	parser csv.Parser
	writer *db.Service
}

func New(config Config, p csv.Parser, service *db.Service) *Application {
	return &Application{
		Config: config,
		parser: p,
		writer: service,
	}
}

func (app *Application) Run(ctx context.Context) {
	fmt.Println("reading from source: ", app.SourcePath)
	sourceFile, err := os.Open(app.SourcePath)
	if err != nil {
		panic(err)
	}
	defer sourceFile.Close()

	pipe := app.writer.Start(ctx)

	if err := app.parser.Consume(sourceFile, pipe); err != nil {
		fmt.Println("Validation error:", err)
	}

}
