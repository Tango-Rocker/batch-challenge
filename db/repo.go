package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type DataRecord struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
}

// repository handles database operations.
type repository struct {
	pool *pgxpool.Pool
	l    *slog.Logger
}

// newRepository creates a new repository with a connection pool.
func newRepository(ctx context.Context, dataSourceName string, l *slog.Logger) (*repository, error) {
	pool, err := pgxpool.New(ctx, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	return &repository{pool: pool, l: l}, nil
}

// InsertData inserts a slice of DataRecord into the database.
func (r *repository) InsertData(ctx context.Context, records []DataRecord) error {
	tx, err := r.pool.Begin(ctx)
	r.l.Info("Transaction started")
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // This will be ignored if tx.Commit() is called

	batch := &pgx.Batch{}
	r.l.Info("Preparing batch size: %d", len(records))
	for _, record := range records {
		batch.Queue("INSERT INTO Transactions (id, value) VALUES ($1, $2)", record.ID, record.Value)
	}
	r.l.Info("Batch created")

	results := tx.SendBatch(ctx, batch)
	defer results.Close()

	for range records {
		if _, err := results.Exec(); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// Close terminates the database connection pool.
func (r *repository) Close() {
	r.pool.Close()
}
