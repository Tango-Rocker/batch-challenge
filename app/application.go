package app

import (
	"context"
	"fmt"
	"github.com/Tango-Rocker/batch-challange/business"
	"github.com/Tango-Rocker/batch-challange/csv"
	"log/slog"
	"os"
	"time"
)

// TODO: i should refactor this a bit more with reader/writer interfaces
// dont really care what kind of io is it really

type Application struct {
	sourcePath string
	parser     csv.Parser
	writer     *business.Writer
	l          *slog.Logger
}

func New(path string, p csv.Parser, w *business.Writer, l *slog.Logger) *Application {
	return &Application{
		sourcePath: path,
		parser:     p,
		writer:     w,
		l:          l,
	}
}

func (app *Application) Run(ctx context.Context) {
	fmt.Println("reading from source: ", app.sourcePath)
	sourceFile, err := os.Open(app.sourcePath)
	if err != nil {
		panic(err)
	}
	defer sourceFile.Close()

	app.writer.Launch(ctx)

	err = app.parser.Consume(newExecutionID(sourceFile.Name()), sourceFile, app.writer)
	if err != nil {
		app.l.Error(err.Error())
	}
}

func newExecutionID(srcName string) string {
	return fmt.Sprintf("batch-%d-%d-%s", os.Getpid(), time.Now().UnixNano(), srcName)
}
