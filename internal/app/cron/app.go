package cron

import (
	"context"
	"plata/internal/clients/exchange"
	"plata/internal/common/log"
	"plata/internal/config"
	"plata/internal/domain/quote"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	qr "plata/internal/repository/quote"
)

type Service struct {
	repo     qr.Repository
	fetcher  exchange.ExternalRateFetcher
	cron     *cron.Cron
	schedule string
	log      log.Logger
}

func New(cfg config.CronConfig, repo qr.Repository, fetcher exchange.ExternalRateFetcher, log log.Logger) *Service {
	return &Service{
		cron:     cron.New(),
		repo:     repo,
		fetcher:  fetcher,
		schedule: cfg.Schedule,
		log:      log,
	}
}

func (s *Service) Run() error {
	s.log.Infof("Cron job registered with schedule: %s", s.schedule)
	_, err := s.cron.AddFunc(s.schedule, func() {
		s.log.Info("Running scheduled quote update...")
		ctx := context.Background()
		quotes, err := s.repo.GetInProgressQuotes(ctx)
		s.log.Infof("Found %d in-progress quotes to update", len(quotes))
		if err != nil {
			s.log.Errorf("Failed to fetch in-progress quotes: %v", err)
			return
		}
		s.updateQuotesBatch(ctx, quotes)
	})
	if err != nil {
		return err
	}
	s.cron.Start()
	return nil
}

func (s *Service) updateQuotesBatch(ctx context.Context, quotes []*quote.Quote) {
	type quoteGroup struct {
		targets []string
		quotes  []*quote.Quote
	}
	grouped := make(map[string]*quoteGroup, len(quotes))
	for _, q := range quotes {
		base, target, ok := parseCurrencyPair(q.Currency)
		if !ok {
			s.log.Warnf("Invalid currency pair format: %s", q.Currency)
			continue
		}
		if _, exists := grouped[base]; !exists {
			grouped[base] = &quoteGroup{}
		}
		grouped[base].targets = append(grouped[base].targets, target)
		grouped[base].quotes = append(grouped[base].quotes, q)
	}

	for base, group := range grouped {
		targets := unique(group.targets)
		rates, err := s.fetcher.FetchRates(ctx, base, targets)
		if err != nil {
			s.log.Errorf("Failed to fetch rates for group base=%s targets=%v: %v", base, targets, err)
			return
		}
		for _, q := range group.quotes {
			_, target, ok := parseCurrencyPair(q.Currency)
			if !ok {
				continue
			}
			rate, ok := rates[target]
			if !ok {
				s.log.Errorf("Missing rate in API response: %s", q.Currency)
				return
			}
			q.Amount = rate
			q.Status = quote.StatusDone
			q.UpdatedAt = time.Now()
			if err = s.repo.Update(ctx, q); err != nil {
				s.log.Errorf("Failed to update quote in DB: %v", err)
				return
			}
			s.log.Infof("Quote updated: id=%s currency=%s amount=%.4f updated_at=%s",
				q.ID, q.Currency, q.Amount, q.UpdatedAt.Format(time.RFC3339),
			)
		}
	}
}

func (s *Service) Stop() {
	s.log.Info("Stopping cron service...")
	if s.cron != nil {
		s.cron.Stop()
	}
}

func parseCurrencyPair(pair string) (base, target string, ok bool) {
	parts := strings.Split(pair, "/")
	if len(parts) != 2 {
		return "", "", false
	}
	return parts[0], parts[1], true
}

func unique(items []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(items))
	for _, item := range items {
		if _, ok := seen[item]; !ok {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
