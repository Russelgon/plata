package quote

import (
	"context"
	"plata/internal/domain/quote"
)

type QuoteRepository interface {
	Save(ctx context.Context, q *quote.Quote) error
	GetByID(ctx context.Context, id string) (*quote.Quote, error)
	GetLatestByCurrency(ctx context.Context, currency string) (*quote.Quote, error)
	Update(ctx context.Context, q *quote.Quote) error
	GetByIdempotencyKey(ctx context.Context, key string) (*quote.Quote, error)
	GetInProgressQuotes(ctx context.Context) ([]*quote.Quote, error)
}

type QuoteClient interface {
	RequestUpdate(ctx context.Context, currency, idemKey string) (string, error)
	GetByID(ctx context.Context, id string) (*quote.Quote, error)
	GetLatestByCurrency(ctx context.Context, currency string) (*quote.Quote, error)
}
