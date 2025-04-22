package test

import (
	"context"
	"plata/internal/common/log"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"plata/internal/domain/quote"
	qs "plata/internal/services/quote"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) GetInProgressQuotes(ctx context.Context) ([]*quote.Quote, error) {
	args := m.Called(ctx)
	if quotes, ok := args.Get(0).([]*quote.Quote); ok {
		return quotes, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRepo) GetByID(ctx context.Context, id string) (*quote.Quote, error) {
	args := m.Called(ctx, id)
	if quote, ok := args.Get(0).(*quote.Quote); ok {
		return quote, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRepo) GetLatestByCurrency(ctx context.Context, currency string) (*quote.Quote, error) {
	args := m.Called(ctx, currency)
	if quote, ok := args.Get(0).(*quote.Quote); ok {
		return quote, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRepo) GetByIdempotencyKey(ctx context.Context, key string) (*quote.Quote, error) {
	args := m.Called(ctx, key)
	if quote, ok := args.Get(0).(*quote.Quote); ok {
		return quote, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRepo) Save(ctx context.Context, q *quote.Quote) error {
	args := m.Called(ctx, q)
	return args.Error(0)
}

func (m *mockRepo) Update(ctx context.Context, q *quote.Quote) error {
	args := m.Called(ctx, q)
	return args.Error(0)
}

type mockFetcher struct {
	mock.Mock
}

func (m *mockFetcher) FetchRates(ctx context.Context, base string, symbols []string) (map[string]float64, error) {
	args := m.Called(ctx, base, symbols)
	if rates, ok := args.Get(0).(map[string]float64); ok {
		return rates, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestRequestUpdate_Success(t *testing.T) {
	repo := new(mockRepo)
	fetcher := new(mockFetcher)

	service := qs.New(repo, fetcher, log.NewZapLogger())

	ctx := context.Background()
	currency := "EUR/USD"
	idemKey := uuid.NewString()

	repo.On("GetByIdempotencyKey", ctx, idemKey).Return((*quote.Quote)(nil), nil)
	repo.On("Save", ctx, mock.AnythingOfType("*quote.Quote")).Return(nil)

	id, err := service.RequestUpdate(ctx, currency, idemKey)

	assert.NoError(t, err)
	assert.NotEmpty(t, id)

	repo.AssertExpectations(t)
}

func TestRequestUpdate_UnsupportedCurrency(t *testing.T) {
	repo := new(mockRepo)
	fetcher := new(mockFetcher)
	service := qs.New(repo, fetcher, log.NewZapLogger())

	ctx := context.Background()
	currency := "GBP/JPY"
	idemKey := uuid.NewString()

	id, err := service.RequestUpdate(ctx, currency, idemKey)

	assert.ErrorIs(t, err, quote.ErrUnsupportedCurrencyPair)
	assert.Empty(t, id)
}
