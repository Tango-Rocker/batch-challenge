package app

import (
	"context"
	"fmt"
	bs "github.com/Tango-Rocker/batch-challenge/business"
	"github.com/Tango-Rocker/batch-challenge/csv"
	"log/slog"
	"os"
	"time"
)

type Application struct {
	sourcePath string
	relay      *bs.StreamRelayService
	writer     *bs.Writer
	summary    *bs.SummaryService
	parser     csv.Parser
	l          *slog.Logger
}

func New(path string, p csv.Parser, w *bs.Writer, r *bs.StreamRelayService, s *bs.SummaryService, l *slog.Logger) *Application {
	return &Application{
		sourcePath: path,
		relay:      r,
		parser:     p,
		writer:     w,
		summary:    s,
		l:          l,
	}
}

// TODO: all services should implement a simple interface getInput/output, singal channel, and launch
func (app *Application) Run(ctx context.Context) {
	//open source file from where we will read the .csv data
	app.l.Info("Opening source file")
	sourceFile, err := os.Open(app.sourcePath)
	if err != nil {
		panic(err)
	}
	defer sourceFile.Close()

	//the relay will propagate the data to the writer and summarizer concurrently
	app.l.Info("Subscribing services to relay")
	app.relay.Subscribe("writer", app.writer.GetInputChannel())
	app.relay.Subscribe("summarizer", app.summary.GetInputChannel())

	//start services
	app.l.Info("Launching services")
	app.writer.Launch(ctx)
	app.summary.Launch(ctx)
	app.relay.Launch(ctx)

	//start parsing
	app.l.Info("Starting parsing")
	err = app.parser.Consume(newExecutionID(sourceFile.Name()), sourceFile, app.relay.GetInputChannel())
	if err != nil {
		app.l.Error(err.Error())
	}

	select {}
}

func newExecutionID(srcName string) string {
	return fmt.Sprintf("batch-%d-%d-%s", os.Getpid(), time.Now().UnixNano(), srcName)
}
