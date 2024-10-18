package controller

import (
	"context"
	"faraway/internal/service"
	"fmt"
	"log"
	"net"
)

// todo add interface

// todo add docs + split on server & client?
type Dialer struct {
	powService service.POWServiceI
}

func NewDialer(powService service.POWServiceI) *Dialer {
	return &Dialer{
		powService: powService,
	}
}

// todo add docs + test + use ctx (add deadline)
func (h Dialer) GetQuote(ctx context.Context, conn net.Conn) ([]byte, error) {
	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			log.Printf("Failed to close conn: %v", closeErr) // todo заменить на zero log or just mark it thought?
		}
	}()

	challenge, err := ReadMessage(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to read challenge: %w", err)
	}

	solution := h.powService.Solve(challenge)
	if err = WriteMessage(conn, solution); err != nil {
		return nil, fmt.Errorf("failed to write solution: %w", err)
	}

	quote, err := ReadMessage(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to read quote: %w", err)
	}

	return quote, nil
}
