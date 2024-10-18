package handler

import (
	"context"
	"faraway/internal/service"
	"fmt"
	"net"
)

type ClientNetI interface {
	GetQuote(ctx context.Context, conn net.Conn) ([]byte, error)
}

// ClientNet wrapper on TCP client.
// Split network handling & service logic.
type ClientNet struct {
	powService service.POWServiceI
}

func NewClientNet(powService service.POWServiceI) *ClientNet {
	return &ClientNet{powService: powService}
}

// GetQuote from a server by conn.
// Before getting a quote we need to solve POW challenge.
func (cn ClientNet) GetQuote(_ context.Context, conn net.Conn) ([]byte, error) {
	challenge, err := ReadMessage(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to read challenge: %w", err)
	}

	solution := cn.powService.Solve(challenge)
	if err = WriteMessage(conn, solution); err != nil {
		return nil, fmt.Errorf("failed to write solution: %w", err)
	}

	quote, err := ReadMessage(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to read quote: %w", err)
	}

	return quote, nil
}
