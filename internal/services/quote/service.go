package quote

import (
	"context"
	"github.com/google/uuid"
	"plata/internal/common/log"

	"plata/internal/clients/exchange"
	"plata/internal/domain/quote"
	"time"
)

type Service struct {
	repo    QuoteRepository
	fetcher exchange.ExternalRateFetcher
	log     log.Logger
}

func New(
	repo QuoteRepository,
	fetcher exchange.ExternalRateFetcher,
	log log.Logger,
) *Service {
	return &Service{
		repo:    repo,
		fetcher: fetcher,
		log:     log,
	}
}

func (s *Service) RequestUpdate(ctx context.Context, currency, idemKey string) (string, error) {
	s.log.Infof("RequestUpdate called with currency: %s, idempotency key: %s", currency, idemKey)

	if _, ok := quote.AllowedPairs[currency]; !ok {
		s.log.Warnf("Unsupported currency pair: %s", currency)
		return "", quote.ErrUnsupportedCurrencyPair
	}

	if idemKey != "" {
		s.log.Infof("Checking idempotency key: %s", idemKey)
		existing, err := s.repo.GetByIdempotencyKey(ctx, idemKey)
		if err != nil {
			s.log.Errorf("Error checking idempotency key: %v", err)
			return "", err
		}
		if existing != nil {
			s.log.Infof("Found existing quote with idempotency key: %s, ID: %s", idemKey, existing.ID)
			return existing.ID, nil
		}
	}

	id := uuid.NewString()
	q := &quote.Quote{
		ID:             id,
		Currency:       currency,
		Status:         quote.StatusInProgress,
		UpdatedAt:      time.Now(),
		IdempotencyKey: idemKey,
	}

	s.log.Infof("Saving new quote: %+v", q)
	if err := s.repo.Save(ctx, q); err != nil {
		s.log.Errorf("Failed to save quote: %v", err)
		return "", err
	}

	s.log.Infof("Quote saved successfully with ID: %s", id)
	return id, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*quote.Quote, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetLatestByCurrency(ctx context.Context, currency string) (*quote.Quote, error) {
	return s.repo.GetLatestByCurrency(ctx, currency)
}
