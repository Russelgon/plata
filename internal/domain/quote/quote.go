package quote

import (
	"encoding/json"
	"time"
)

type Quote struct {
	ID             string    `json:"id"`
	Currency       string    `json:"currency"`
	Amount         float64   `json:"amount"`
	UpdatedAt      time.Time `json:"updated_at"`
	Status         Status    `json:"status"`
	IdempotencyKey string    `json:"idempotency_key,omitempty"`
}

var AllowedPairs = map[string]struct{}{
	"EUR/USD": {},
	"EUR/MXN": {},
	"EUR/RUB": {},
}

type Status int

const (
	StatusUnspecified Status = iota
	StatusInProgress
	StatusDone
)

func ToString(s Status) string {
	switch s {
	case StatusInProgress:
		return "in_progress"
	case StatusDone:
		return "done"
	default:
		return "unspecified"
	}
}

func FromString(s string) Status {
	switch s {
	case "in_progress":
		return StatusInProgress
	case "done":
		return StatusDone
	default:
		return StatusUnspecified
	}
}

func (q Quote) MarshalJSON() ([]byte, error) {
	type Alias Quote
	return json.Marshal(&struct {
		Status string `json:"status"`
		*Alias
	}{
		Status: ToString(q.Status),
		Alias:  (*Alias)(&q),
	})
}
