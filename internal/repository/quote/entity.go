package quote

import (
	"database/sql"
	"time"
)

type entity struct {
	ID             string         `db:"id"`
	Currency       string         `db:"currency"`
	Amount         float64        `db:"amount"`
	UpdatedAt      time.Time      `db:"updated_at"`
	Status         string         `db:"status"`
	IdempotencyKey sql.NullString `db:"idempotency_key"`
}
