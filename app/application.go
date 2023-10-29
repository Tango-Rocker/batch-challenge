package app

import (
	"context"
	"fmt"
	"github.com/Tango-Rocker/batch-challange/bussines"
	"github.com/Tango-Rocker/batch-challange/csv"
	"os"
	"time"
)

// TODO: i should refactor this a bit more with reader/writer interfaces
// dont really care what kind of io is it really

type Application struct {
	*Config
	parser csv.Parser
	writer *bussines.Worker
}

func New(config *Config, p csv.Parser, w *bussines.Worker) *Application {
	return &Application{
		Config: config,
		parser: p,
		writer: w,
	}
}

func (app *Application) Run(ctx context.Context) {
	fmt.Println("reading from source: ", app.SourcePath)
	sourceFile, err := os.Open(app.SourcePath)
	if err != nil {
		panic(err)
	}
	defer sourceFile.Close()

	go app.writer.Start(ctx)

	if err := app.parser.Consume(newExecutionID(sourceFile.Name()), sourceFile, app.writer); err != nil {
		fmt.Println("Validation error:", err)
	}

}

func newExecutionID(srcName string) string {
	return fmt.Sprintf("batch-%d-%d-%s", os.Getpid(), time.Now().UnixNano(), srcName)
}
