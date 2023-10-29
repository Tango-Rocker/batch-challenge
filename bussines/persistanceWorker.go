package bussines

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Tango-Rocker/batch-challange/csv"
	"github.com/Tango-Rocker/batch-challange/db"
	"github.com/Tango-Rocker/batch-challange/model"
	"log"
	"log/slog"
	"time"
)

// Worker processes data and communicates with the repository.
type Worker struct {
	running      bool
	repo         *db.Repository
	dataChannel  chan []byte
	bufferSize   int
	buffer       []model.Transaction
	flushTimeout time.Duration
	l            *slog.Logger
}

func (w *Worker) Write(p []byte) (n int, err error) {
	w.dataChannel <- p
	return len(p), nil
}

// NewWorker creates a new Worker instance.
func NewWorker(cfg *Config, repo *db.Repository, l *slog.Logger) *Worker {
	return &Worker{
		repo:         repo,
		dataChannel:  make(chan []byte),
		bufferSize:   cfg.BufferSize,
		buffer:       make([]model.Transaction, 0, cfg.BufferSize),
		flushTimeout: time.Duration(cfg.FlushTimeout) * time.Millisecond,
		l:            l,
	}
}

// Start begins processing data from the data channel.
func (w *Worker) Start(ctx context.Context) {
	if w.running {
		return
	}
	w.l.Info("Starting BufferedInsert worker")
	w.running = true
	timeout := time.NewTimer(w.flushTimeout)

	for {
		select {
		case <-ctx.Done():
			w.flushBuffer(ctx)
			w.running = false
			w.l.Info("Stopping BufferedInsert worker")

			return
		case jsonData := <-w.dataChannel:
			//TODO: dropped data should be logged and stored in the database for later processing
			var record csv.Record
			if err := json.Unmarshal(jsonData, &record); err != nil {
				w.l.Error("Error unmarshalling data: %v", err)
				w.l.Error(string(jsonData))
				continue
			}

			trx, err := mapRecordToTransaction(&record)
			if err != nil {
				w.l.Error("Error mapping record to transaction: %v", err)
				w.l.Error(string(jsonData))
				continue
			}

			w.buffer = append(w.buffer, *trx)
			if len(w.buffer) >= w.bufferSize {
				w.flushBuffer(ctx)
				timeout.Reset(w.flushTimeout)
			}
		case <-timeout.C:
			if len(w.buffer) > 0 {
				w.flushBuffer(ctx)
			}
			timeout.Reset(w.flushTimeout)
		}
	}
}

// flushBuffer inserts the buffered data records into the database.
func (w *Worker) flushBuffer(ctx context.Context) {
	if err := w.repo.InsertData(ctx, w.buffer); err != nil {
		log.Printf("Error inserting data: %v", err)
	}
	w.buffer = w.buffer[:0] // Clear the buffer
}
func mapRecordToTransaction(record *csv.Record) (*model.Transaction, error) {
	valuesJSON, err := json.Marshal(record.Values)
	if err != nil {
		return nil, fmt.Errorf("error marshaling values: %v", err)
	}

	var transaction model.Transaction
	err = json.Unmarshal(valuesJSON, &transaction)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling into Transaction: %v", err)
	}

	transaction.BatchId = record.ExecutionId
	transaction.Line = int(record.Line)     // Converting uint64 to int
	transaction.Offset = int(record.Offset) // Converting uint64 to int

	return &transaction, nil
}
