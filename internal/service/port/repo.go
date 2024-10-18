package port

import (
	"context"
	"faraway/internal/domain"
)

type QuoteRepoI interface {
	GetRand(ctx context.Context) domain.Quote
}
