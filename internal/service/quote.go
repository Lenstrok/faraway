package service

import (
	"context"
	"faraway/internal/domain"
	"faraway/internal/service/port"
)

type QuoteServiceI interface {
	GetRand(ctx context.Context) domain.Quote
}

// todo docs
type QuoteService struct {
	quoteRepo port.QuoteRepoI
}

func NewQuoteService(quoteRepo port.QuoteRepoI) *QuoteService {
	return &QuoteService{quoteRepo: quoteRepo}
}

// todo + add ctx may use in DB
func (qs QuoteService) GetRand(ctx context.Context) domain.Quote {
	return qs.quoteRepo.GetRand(ctx)
}
