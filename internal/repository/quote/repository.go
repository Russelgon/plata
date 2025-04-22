package quote

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"plata/internal/domain/quote"
)

type Repository struct {
	dbP *sqlx.DB
	dbR *sqlx.DB
}

func New(primary, replica *sqlx.DB) *Repository {
	return &Repository{
		dbP: primary,
		dbR: replica,
	}
}

func (r *Repository) GetByID(ctx context.Context, id string) (*quote.Quote, error) {
	const query = `
		SELECT id, currency, amount, status, updated_at, idempotency_key
		FROM quotes
		WHERE id = $1
	`
	var e entity
	err := r.dbR.GetContext(ctx, &e, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, quote.ErrQuoteNotFound
		}
		return nil, fmt.Errorf("failed to get quote by ID: %w", err)
	}
	return toDomain(&e), nil
}

func (r *Repository) GetLatestByCurrency(ctx context.Context, currency string) (*quote.Quote, error) {
	query := `
		SELECT id, currency, amount, status, updated_at, idempotency_key
		FROM quotes
		WHERE currency = $1
		ORDER BY updated_at DESC
		LIMIT 1
	`
	var e entity
	err := r.dbR.GetContext(ctx, &e, query, currency)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, quote.ErrQuoteNotFound
		}
		return nil, fmt.Errorf("failed to get latest quote: %w", err)
	}

	return toDomain(&e), nil
}

func (r *Repository) Update(ctx context.Context, q *quote.Quote) error {
	_, err := r.dbP.ExecContext(ctx,
		`UPDATE quotes 
		 SET amount = $1, status = $2, updated_at = $3 
		 WHERE id = $4`,
		q.Amount, quote.ToString(q.Status), q.UpdatedAt, q.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update quote: %w", err)
	}
	return nil
}

func (r *Repository) Save(ctx context.Context, q *quote.Quote) error {
	query := `
		INSERT INTO quotes (id, currency, amount, status, updated_at, idempotency_key)
		VALUES (:id, :currency, :amount, :status, :updated_at, :idempotency_key)
	`
	rec := toEntity(q)
	_, err := r.dbP.NamedExecContext(ctx, query, rec)
	if err != nil {
		return fmt.Errorf("failed to save quote: %w", err)
	}
	return nil
}

func (r *Repository) GetByIdempotencyKey(ctx context.Context, key string) (*quote.Quote, error) {
	query := `
		SELECT id, currency, amount, status, updated_at, idempotency_key
		FROM quotes
		WHERE idempotency_key = $1
	`

	var e entity
	err := r.dbR.GetContext(ctx, &e, query, key)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get quote by idempotency key: %w", err)
	}
	return toDomain(&e), nil
}

func (r *Repository) GetInProgressQuotes(ctx context.Context) ([]*quote.Quote, error) {
	const query = `
		SELECT id, currency, amount, status, updated_at, idempotency_key
		FROM quotes
		WHERE status = $1
	`
	var e []entity
	err := r.dbR.SelectContext(ctx, &e, query, quote.ToString(quote.StatusInProgress))
	if err != nil {
		return nil, fmt.Errorf("failed to get quotes: %w", err)
	}

	quotes := make([]*quote.Quote, len(e))
	for i := range e {
		quotes[i] = toDomain(&e[i])
	}
	return quotes, nil
}
