package db

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

type Service struct {
	repo   *repository
	worker *Worker
}

const postgresConnectionFormat = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable"

func NewService(cfg Config, ctx context.Context, l *slog.Logger) (*Service, error) {

	url := fmt.Sprintf(postgresConnectionFormat,
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Pass,
		cfg.Name)

	repo, err := newRepository(ctx, url, l)
	if err != nil {
		return nil, err
	}

	worker := NewWorker(repo, cfg.BulkSize, time.Second*time.Duration(cfg.FlushTimeout))

	return &Service{
		repo:   repo,
		worker: worker,
	}, nil
}

func (s *Service) Start(ctx context.Context) chan []byte {
	go s.worker.Start(ctx)
	return s.worker.dataChannel
}
