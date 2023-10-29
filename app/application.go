package app

import (
	"fmt"
	"github.com/Tango-Rocker/batch-challange/csv"
	"io"
	"os"
	"sync"
)

type Application struct {
	Config
	parser csv.Parser
}

func New(config Config, p csv.Parser) *Application {
	return &Application{
		Config: config,
		parser: p,
	}
}

func (app *Application) Run() {
	fmt.Println("reading from source: ", app.SourcePath)
	sourceFile, err := os.Open(app.SourcePath)
	if err != nil {
		panic(err)
	}
	defer sourceFile.Close()

	recordPipe := make(chan []string)
	var wg sync.WaitGroup

	// Start Schema reading and validation in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := app.parser.Consume(sourceFile, recordPipe); err != nil {
			fmt.Println("Validation error:", err)
			close(recordPipe)
		}
	}()

	// Start processing the recordPipe in a separate goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		processRecords(recordPipe, os.Stdout)
	}()

	wg.Wait()
}

func processRecords(records <-chan []string, output io.Writer) {
	i := 0
	for record := range records {
		// Process the record to construct the desired payload
		// For the sake of the example, we just print it out.
		_, err := fmt.Fprintln(output, record)
		if err != nil {
			fmt.Printf("Error writing record %d: %s\n", i, err.Error())
		}
		i++
	}
	fmt.Printf("Processed %d records\n", i)
}
