package repo

import (
	"context"
	"encoding/json"
	"faraway/internal/domain"
	"math/rand"
)

type QuoteRepo struct {
	quotes []domain.Quote
}

// TODO add doc + test + describe json file fmt
func NewQuoteRepo(quotesBytes []byte) (*QuoteRepo, error) {
	var book struct {
		Quotes []domain.Quote `json:"quotes"`
	}

	if err := json.Unmarshal(quotesBytes, &book); err != nil {
		return nil, err
	}

	return &QuoteRepo{quotes: book.Quotes}, nil
}

// TODO add doc + test
func (qr QuoteRepo) GetRand(_ context.Context) domain.Quote {
	return qr.quotes[rand.Intn(len(qr.quotes))]
}
