package db

import (
	"context"
	"database/sql"
)

type DataRecord struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
}

// repository handles database operations.
type repository struct {
	db *sql.DB
}

// newRepository creates a new repository.
func newRepository(dataSourceName string) (*repository, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &repository{db: db}, nil
}

// InsertData inserts a slice of DataRecord into the database.
func (r *repository) InsertData(ctx context.Context, records []DataRecord) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO Transactions (id, value) VALUES ($1, $2)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, record := range records {
		if _, err := stmt.ExecContext(ctx, record.ID, record.Value); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

// Close terminates the database connection.
func (r *repository) Close() error {
	return r.db.Close()
}
