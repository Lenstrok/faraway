package infra

import (
	"context"
	"errors"
	"faraway/internal/adapter/handler"
	"faraway/internal/config"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type Listener struct {
	handler  handler.ServerNetI
	listener net.Listener

	wg       *sync.WaitGroup
	cfg      config.Server
	cancelFn context.CancelFunc
}

// NewListener with tcp network.
func NewListener(handler handler.ServerNetI, cfg config.Server) (*Listener, error) {
	const network = "tcp"

	listener, err := net.Listen(network, ":"+cfg.Port)
	if err != nil {
		return nil, fmt.Errorf("failed to get net listener: %w", err)
	}

	return &Listener{
		handler:  handler,
		listener: listener,
		cfg:      cfg,
		wg:       &sync.WaitGroup{},
	}, nil
}

// Start to accept requests & handle them concurrently.
// May be stopped by calling Stop method.
func (l *Listener) Start(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	l.cancelFn = cancel

	l.wg.Add(1)

	go func() {
		defer l.wg.Done()

		for {
			conn, err := l.listener.Accept()
			if errors.Is(err, net.ErrClosed) {
				log.Printf("Listener closed.")
				return
			} else if err != nil {
				continue
			}

			log.Printf("Listener accept.")

			_ = conn.SetDeadline(time.Now().Add(l.cfg.ConnTimeout))

			l.wg.Add(1)

			go func() {
				defer l.wg.Done()
				defer func() {
					if closeErr := conn.Close(); closeErr != nil {
						log.Printf("Failed to close conn: %v", closeErr)
					}
				}()

				if handleErr := l.handler.HandleQuoteReq(ctx, conn); handleErr != nil {
					log.Printf("Failed to handle quote request: %v", handleErr)
				}
			}()
		}
	}()

	<-ctx.Done()

	if err := l.listener.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
		log.Printf("Failed to close listener: %v", err)
	}
}

// Stop the listener.
// Method is blocked until all goroutines stopped.
func (l *Listener) Stop() {
	l.cancelFn()

	l.wg.Wait()
}
