package api

type UpdateQuoteRequest struct {
	Currency string `json:"currency" binding:"required,len=7"`
}

type UpdateQuoteResponse struct {
	UpdateID string `json:"update_id"`
}
