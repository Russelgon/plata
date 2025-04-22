package api

import (
	"errors"
	"github.com/google/uuid"
	"net/http"

	"github.com/gin-gonic/gin"
	dq "plata/internal/domain/quote"
	"plata/internal/services/quote"
)

type Handler struct {
	QuoteService quote.QuoteClient
}

func NewHandler(quoteService quote.QuoteClient) *Handler {
	return &Handler{QuoteService: quoteService}
}

type ErrorResponse struct {
	Error   string      `json:"error"`
	Status  int         `json:"code,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

type SuccessResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// UpdateQuote updates a currency pair quote
// @Summary Update a quote
// @Description Asynchronously request a quote update for a currency pair
// @Tags quotes
// @Accept json
// @Produce json
// @Param request body UpdateQuoteRequest true "Currency pair"
// @Success 200 {object} SuccessResponse{data=UpdateQuoteResponse}
// @Failure 400,500 {object} ErrorResponse "Error response"
// @Router /quotes/update [post]
func (h *Handler) UpdateQuote(c *gin.Context) {
	var req UpdateQuoteRequest
	idemKey := c.GetHeader("Idempotency-Key")
	if idemKey == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Idempotency-Key is required",
		})
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "invalid request",
			Details: err.Error(),
		})
		return
	}
	id, err := h.QuoteService.RequestUpdate(c.Request.Context(), req.Currency, idemKey)
	if err != nil {
		if errors.Is(err, dq.ErrUnsupportedCurrencyPair) {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Status:  http.StatusBadRequest,
				Error:   "currency pair not supported currently",
				Details: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Status:  http.StatusInternalServerError,
			Error:   "failed to request update",
			Details: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{
		Status:  http.StatusOK,
		Message: "quote update requested",
		Data:    UpdateQuoteResponse{UpdateID: id},
	})
}

// GetQuoteByID retrieves a quote by update ID
// @Summary Retrieve quote by ID
// @Tags quotes
// @Produce json
// @Param id path string true "Quote update ID"
// @Success 200 {object} SuccessResponse{data=dq.Quote}
// @Failure 400,404,500 {object} ErrorResponse "Error response"
// @Router /quotes/{id} [get]
func (h *Handler) GetQuoteByID(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "invalid ID format, must be a UUID",
			Details: err.Error(),
		})
		return
	}
	q, err := h.QuoteService.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, dq.ErrQuoteNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Status:  http.StatusNotFound,
				Error:   "unable to find quote with such ID",
				Details: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Status:  http.StatusInternalServerError,
			Error:   "internal error",
			Details: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{
		Status:  http.StatusOK,
		Message: "quote found",
		Data:    q,
	})
}

// GetLatestQuote returns the latest quote for a given currency pair
// @Summary Get latest quote
// @Description Returns the most recent quote for a currency pair (e.g. EUR/USD)
// @Tags quotes
// @Accept json
// @Produce json
// @Param currency query string true "Currency pair" Enums(EUR/USD, EUR/MXN, EUR/RUB)
// @Success 200 {object} SuccessResponse{data=dq.Quote}
// @Failure 400,404,500 {object} ErrorResponse "Error response"
// @Router /quotes/latest [get]
func (h *Handler) GetLatestQuote(c *gin.Context) {
	currency := c.Query("currency")
	if currency == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "currency query parameter is required",
			Details: dq.ErrUnsupportedCurrencyPair,
		})
		return
	}
	q, err := h.QuoteService.GetLatestByCurrency(c.Request.Context(), currency)
	if err != nil {
		if errors.Is(err, dq.ErrQuoteNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Status:  http.StatusNotFound,
				Error:   "quote not found",
				Details: err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Status:  http.StatusInternalServerError,
				Error:   "internal error",
				Details: err.Error(),
			})
			return
		}
	}
	c.JSON(http.StatusOK, SuccessResponse{
		Status:  http.StatusOK,
		Message: "latest quote retrieved",
		Data:    q,
	})
}
