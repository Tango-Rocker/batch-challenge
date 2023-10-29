package db

import (
	"context"
	"encoding/json"
	"log"
	"time"
)

// Worker processes data and communicates with the repository.
type Worker struct {
	repo         *Repository
	dataChannel  chan []byte
	bufferSize   int
	buffer       []DataRecord
	flushTimeout time.Duration
}

// NewWorker creates a new Worker instance.
func NewWorker(repo *Repository, bufferSize int, flushTimeout time.Duration) *Worker {
	return &Worker{
		repo:         repo,
		dataChannel:  make(chan []byte),
		bufferSize:   bufferSize,
		buffer:       make([]DataRecord, 0, bufferSize),
		flushTimeout: flushTimeout,
	}
}

// Start begins processing data from the data channel.
func (w *Worker) Start(ctx context.Context) {
	timeout := time.NewTimer(w.flushTimeout)

	for {
		select {
		case <-ctx.Done():
			return
		case jsonData := <-w.dataChannel:
			var record DataRecord
			if err := json.Unmarshal(jsonData, &record); err != nil {
				log.Printf("Error unmarshalling JSON: %v", err)
				continue
			}
			w.buffer = append(w.buffer, record)
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

func example() {
	ctx := context.Background()

	repo, err := NewRepository("")
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()

	worker := NewWorker(repo, 100, 30*time.Second)
	go worker.Start(ctx)
}
