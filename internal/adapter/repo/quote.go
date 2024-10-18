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

// NewQuoteRepo filled with collection.
// Data in JSON Format: `{"quotes": [{"text": "magic", "author": "mage"}, ..]}`
func NewQuoteRepo(quotesBytes []byte) (*QuoteRepo, error) {
	var book struct {
		Quotes []domain.Quote `json:"quotes"`
	}

	if err := json.Unmarshal(quotesBytes, &book); err != nil {
		return nil, err
	}

	return &QuoteRepo{quotes: book.Quotes}, nil
}

// GetRand quote from the collection.
// Context may be useful to handling DB requests.
func (qr QuoteRepo) GetRand(_ context.Context) domain.Quote {
	return qr.quotes[rand.Intn(len(qr.quotes))]
}
