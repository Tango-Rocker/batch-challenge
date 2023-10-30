package business

import (
	"context"
	"log/slog"
	"sync"
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
		l:             l.With(slog.String("service", "stream-relay")),
		running:       false,
		subscriptions: make(map[string]chan []byte),
		input:         make(chan []byte, 1000),
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
	wg := sync.WaitGroup{}

	for {
		select {
		case <-ctx.Done():
			{
				srv.l.Info("StreamRelayService received done signal")
				srv.running = false
				return
			}
		case data, ok := <-srv.input:
			{
				if !ok {
					srv.l.Info("StreamRelayService input channel closed")
					srv.running = false

					// wait for all the routines to finish before closing the subscriptions
					wg.Wait()
					for k := range srv.subscriptions {
						close(srv.subscriptions[k])
					}
					return
				}

				for k := range srv.subscriptions {
					key := k // otherwise the closure will use the last value of k

					//FIXME: we should use a pool of routines and a c-breaker
					// to prevent an overflow in the system
					wg.Add(1)
					go func() {
						timeout := time.After(5 * time.Second)

						select {
						case srv.subscriptions[key] <- data:
						case <-timeout:
							srv.l.Error("StreamRelayService timed out on subscriptions: " + key)
						}
						wg.Done()
					}()
				}
			}
		}
	}
}
