package business

import (
	"context"
	"encoding/json"
	"github.com/Tango-Rocker/batch-challange/csv"
	"github.com/Tango-Rocker/batch-challange/db"
	"github.com/Tango-Rocker/batch-challange/model"
	"log"
	"log/slog"
	"time"
)

// Writer processes data and communicates with the repository.
type Writer struct {
	running      bool
	repo         *db.Repository
	dataChannel  chan []byte
	bufferSize   int
	buffer       []model.Transaction
	flushTimeout time.Duration
	l            *slog.Logger
}

func (w *Writer) Write(p []byte) (n int, err error) {
	w.dataChannel <- p
	return len(p), nil
}

// NewWriter creates a new Writer instance.
func NewWriter(cfg *WriterConfig, repo *db.Repository, l *slog.Logger) *Writer {
	return &Writer{
		repo:         repo,
		dataChannel:  make(chan []byte, 1000),
		bufferSize:   cfg.BufferSize,
		buffer:       make([]model.Transaction, 0, cfg.BufferSize),
		flushTimeout: time.Duration(cfg.FlushTimeout) * time.Millisecond,
		l:            l.With(slog.String("service", "db-writer")),
	}
}

// Launch begins processing data from the data channel.
func (w *Writer) Launch(ctx context.Context) {
	if !w.running {
		go w.run(ctx)
	}
}

func (w *Writer) run(ctx context.Context) {
	w.l.Info("Starting BatchInsert Service")
	w.running = true
	timeout := time.NewTimer(w.flushTimeout)

	for {
		select {
		case <-ctx.Done():
			w.flushBuffer(ctx)
			w.running = false
			w.l.Info("Stopping BufferedInsert worker")

			return
		case jsonData, ok := <-w.dataChannel:

			if !ok {
				w.l.Info("Data channel closed, stopping service")
				w.flushBuffer(ctx)
				return
			}

			//TODO: dropped data should be logged and stored in the database for later processing
			var record csv.Record
			if err := json.Unmarshal(jsonData, &record); err != nil {
				w.l.Error("Error unmarshalling data: %v", err)
				w.l.Error(string(jsonData))
				//TODO: we need to signal that an error occurred
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
func (w *Writer) flushBuffer(ctx context.Context) {
	if err := w.repo.InsertData(ctx, w.buffer); err != nil {
		log.Printf("Error inserting data: %v", err)
	}
	w.buffer = w.buffer[:0] // Clear the buffer
}

func (w *Writer) GetInputChannel() chan []byte {
	return w.dataChannel
}
