package exchange

import "context"

type ExternalRateFetcher interface {
	FetchRates(ctx context.Context, base string, symbols []string) (map[string]float64, error)
}
