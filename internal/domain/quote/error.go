package quote

import "errors"

var (
	ErrUnsupportedCurrencyPair = errors.New("unsupported currency pair")
	ErrQuoteNotFound           = errors.New("quote not found")
)
