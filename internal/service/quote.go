package service

import (
	"context"
	"faraway/internal/domain"
	"faraway/internal/service/port"
)

type QuoteServiceI interface {
	GetRand(ctx context.Context) domain.Quote
}

type QuoteService struct {
	quoteRepo port.QuoteRepoI
}

func NewQuoteService(quoteRepo port.QuoteRepoI) *QuoteService {
	return &QuoteService{quoteRepo: quoteRepo}
}

// GetRand quote from the collection.
func (qs QuoteService) GetRand(ctx context.Context) domain.Quote {
	return qs.quoteRepo.GetRand(ctx)
}
