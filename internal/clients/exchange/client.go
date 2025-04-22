package exchange

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"plata/internal/common/log"
	"plata/internal/config"
	"strings"
)

type Service struct {
	URL    string
	APIKey string
	client *http.Client
	log    log.Logger
}

func New(cfg config.ExchangeConfig, log log.Logger) *Service {
	return &Service{
		APIKey: cfg.APIKey,
		URL:    cfg.URL,
		client: &http.Client{
			Timeout: cfg.Timeout,
		},
		log: log,
	}
}

func (s *Service) FetchRates(ctx context.Context, base string, targets []string) (map[string]float64, error) {
	urlF, err := buildURL(s.URL, map[string]string{
		"access_key": s.APIKey,
		"base":       base,
		"symbols":    strings.Join(targets, ","),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to parse URL: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlF, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %v", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.log.Errorf("API error: status=%s code=%d", resp.Status, resp.StatusCode)
		return nil, fmt.Errorf("unexpected API error: %s", resp.Status)
	}

	var result struct {
		Success bool               `json:"success"`
		Rates   map[string]float64 `json:"rates"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	if !result.Success {
		s.log.Error("API responded with success = false")
		return nil, fmt.Errorf("API response was not successful")
	}
	return result.Rates, nil
}

func buildURL(base string, params map[string]string) (string, error) {
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	query := u.Query()
	for key, value := range params {
		query.Set(key, value)
	}
	u.RawQuery = query.Encode()
	return u.String(), nil
}
