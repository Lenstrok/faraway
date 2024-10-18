package handler

import (
	"context"
	"faraway/internal/service"
	"fmt"
	"log"
	"net"
)

type ServerNetI interface {
	HandleQuoteReq(ctx context.Context, conn net.Conn) error
}

// todo add docs + split on server & client?
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

// todo add docs + test + use ctx (add deadline)
func (sn ServerNet) HandleQuoteReq(ctx context.Context, conn net.Conn) error {
	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			log.Printf("Failed to close conn: %v", closeErr) // todo заменить на zero log or just mark it thought?
		}
	}()

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
