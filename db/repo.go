package db

import (
	"context"
	"fmt"
	"gorm.io/gorm/logger"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/Tango-Rocker/batch-challenge/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Repository handles database operations.
type Repository struct {
	db *gorm.DB
	l  *slog.Logger
}

// NewRepository creates a new Repository with a GORM connection.
func NewRepository(cfg *Config, l *slog.Logger) (*Repository, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.Host,
		cfg.User,
		cfg.Pass,
		cfg.Name,
		cfg.Port)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		l.Error("failed to connect to database: %v", err)
		return nil, err
	}

	err = db.AutoMigrate(&model.Transaction{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		l.Error("failed to get database connection handle: %v", err)
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(0)

	return &Repository{db: db, l: l}, nil
}

// InsertData inserts a slice of DataRecord into the database using GORM.
func (r *Repository) InsertData(ctx context.Context, records []model.Transaction) error {
	r.l.Info("Inserting data")
	// GORM handles transactions by default, so you just need to use Create method.
	result := r.db.WithContext(ctx).Create(&records)
	if result.Error != nil {
		return result.Error
	}

	r.l.Info(fmt.Sprintf("Inserted %d records", len(records)))
	return nil
}

// Close terminates the database connection pool.
func (r *Repository) Close() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection handle: %v", err)
	}
	sqlDB.Close()
	return nil
}
