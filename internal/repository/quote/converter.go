package quote

import (
	"database/sql"
	"plata/internal/domain/quote"
)

func toDomain(qr *entity) *quote.Quote {
	var key string
	if qr.IdempotencyKey.Valid {
		key = qr.IdempotencyKey.String
	}

	return &quote.Quote{
		ID:             qr.ID,
		Currency:       qr.Currency,
		Amount:         qr.Amount,
		UpdatedAt:      qr.UpdatedAt,
		Status:         quote.FromString(qr.Status),
		IdempotencyKey: key,
	}
}

func toEntity(q *quote.Quote) *entity {
	var key sql.NullString
	if q.IdempotencyKey != "" {
		key = sql.NullString{String: q.IdempotencyKey, Valid: true}
	}
	return &entity{
		ID:             q.ID,
		Currency:       q.Currency,
		Amount:         q.Amount,
		UpdatedAt:      q.UpdatedAt,
		Status:         quote.ToString(q.Status),
		IdempotencyKey: key,
	}
}
