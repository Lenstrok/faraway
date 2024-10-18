package infra

import (
	"context"
	"errors"
	"faraway/internal/adapter/handler"
	"faraway/internal/config"
	"fmt"
	"log"
	"net"
)

type Dialer struct {
	handler handler.ClientNetI

	dialer net.Dialer
	cfg    config.Server
}

func NewDialer(handler handler.ClientNetI, cfg config.Server) *Dialer {
	return &Dialer{
		handler: handler,
		dialer:  net.Dialer{},
		cfg:     cfg,
	}
}

// GetQuote from a server by tcp network.
func (d Dialer) GetQuote(ctx context.Context) (string, error) {
	const network = "tcp"

	conn, err := d.dialer.DialContext(ctx, network, ":"+d.cfg.Port)
	if err != nil {
		return "", fmt.Errorf("failed to connect to server: %w", err)
	}

	defer func() {
		if closeErr := conn.Close(); closeErr != nil && !errors.Is(closeErr, net.ErrClosed) {
			log.Printf("Failed to close dialer: %v", closeErr)
		}
	}()

	quote, err := d.handler.GetQuote(ctx, conn)
	if err != nil {
		return "", fmt.Errorf("failed to get a quote: %w", err)
	}

	return string(quote), nil
}
