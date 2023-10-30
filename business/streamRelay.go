package business

import (
	"context"
	"log/slog"
	"time"
)

type StreamRelayService struct {
	l             *slog.Logger
	running       bool
	subscriptions map[string]chan []byte
	input         chan []byte
}

func NewStreamRelayService(l *slog.Logger) *StreamRelayService {
	return &StreamRelayService{
		l:             l,
		running:       false,
		subscriptions: make(map[string]chan []byte),
		input:         make(chan []byte),
	}
}

func (srv *StreamRelayService) GetInputChannel() chan []byte {
	return srv.input
}

func (srv *StreamRelayService) Subscribe(name string, sub chan []byte) {
	srv.subscriptions[name] = sub
}

func (srv *StreamRelayService) Launch(ctx context.Context) {
	if !srv.running {
		go srv.run(ctx)
	}
}

func (srv *StreamRelayService) run(ctx context.Context) {
	srv.running = true
	srv.l.Info("Starting StreamRelayService")
	for {
		select {
		case <-ctx.Done():
			{
				srv.l.Info("StreamRelayService received done signal")
				srv.running = false
				return
			}
		case data := <-srv.input:
			{
				for k := range srv.subscriptions {
					key := k // otherwise the closure will use the last value of k

					//FIXME: we should use a pool of routines and a c-breaker
					// to prevent an overflow in the system
					go func() {
						timeout := time.After(1 * time.Second)

						select {
						case srv.subscriptions[key] <- data:
						case <-timeout:
							srv.l.Error("StreamRelayService timed out on subscriptions: " + key)
						}
					}()
				}
			}
		}
	}
}
