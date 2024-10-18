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
type Handler struct {
	powService   service.POWServiceI
	quoteService service.QuoteServiceI
}

func NewHandler(powService service.POWServiceI, quoteService service.QuoteServiceI) *Handler {
	return &Handler{
		powService:   powService,
		quoteService: quoteService,
	}
}

// todo add docs + test + use ctx (add deadline)
func (h Handler) HandleQuoteReq(ctx context.Context, conn net.Conn) error {
	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			log.Printf("Failed to close conn: %v", closeErr) // todo заменить на zero log or just mark it thought?
		}
	}()

	challenge := h.powService.Challenge()
	if err := WriteMessage(conn, challenge); err != nil {
		return fmt.Errorf("failed to write challenge: %w", err)
	}

	solution, err := ReadMessage(conn)
	if err != nil {
		return fmt.Errorf("failed to read solution: %w", err)
	}

	if err = h.powService.Verify(challenge, solution); err != nil {
		return fmt.Errorf("failed to verify: %w", err)
	}

	if err = WriteMessage(conn, h.quoteService.GetRand(ctx).Bytes()); err != nil {
		return fmt.Errorf("failed to write quote: %w", err)
	}

	return nil
}
