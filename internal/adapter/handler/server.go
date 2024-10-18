package handler

import (
	"context"
	"faraway/internal/service"
	"fmt"
	"net"
)

type ServerNetI interface {
	HandleQuoteReq(ctx context.Context, conn net.Conn) error
}

// ServerNet wrapper on TCP server.
// Split network handling & service logic.
type ServerNet struct {
	powService   service.POWServiceI
	quoteService service.QuoteServiceI
}

func NewServerNet(powService service.POWServiceI, quoteService service.QuoteServiceI) *ServerNet {
	return &ServerNet{
		powService:   powService,
		quoteService: quoteService,
	}
}

// HandleQuoteReq to the server by conn.
// Before sending a quote we need to send POW challenge and get the right answer.
func (sn ServerNet) HandleQuoteReq(ctx context.Context, conn net.Conn) error {
	challenge := sn.powService.Challenge()
	if err := WriteMessage(conn, challenge); err != nil {
		return fmt.Errorf("failed to write challenge: %w", err)
	}

	solution, err := ReadMessage(conn)
	if err != nil {
		return fmt.Errorf("failed to read solution: %w", err)
	}

	if err = sn.powService.Verify(challenge, solution); err != nil {
		return fmt.Errorf("failed to verify: %w", err)
	}

	quote := sn.quoteService.GetRand(ctx).Bytes()
	if err = WriteMessage(conn, quote); err != nil {
		return fmt.Errorf("failed to write quote: %w", err)
	}

	return nil
}
