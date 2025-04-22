package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"plata/internal/common/log"

	"net/http"
)

type Service struct {
	srv *http.Server
	log log.Logger
}

func NewServer(log log.Logger) *Service {
	return &Service{
		log: log,
	}
}

func (s *Service) InitServer(handler *Handler) error {
	r := gin.Default()
	api := r.Group("/api/v1/quotes")
	{
		// 1. Update quote (POST /api/v1/quotes/update)
		api.POST("/update", handler.UpdateQuote)

		// 2. get quote by ID (GET /api/v1/quotes/:id)
		api.GET("/:id", handler.GetQuoteByID)

		// 3. get last quote by pair (GET /api/v1/quotes/latest/:pair)
		api.GET("/latest", handler.GetLatestQuote)
	}
	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Healthcheck
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	s.srv = &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	if err := s.run(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

func (s *Service) run() error {
	s.log.Info("Starting server on :8080")
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Errorf("HTTP server error: %v", err)
		}
	}()
	return nil
}

func (s *Service) Stop(ctx context.Context) {
	s.log.Info("Shutting down server...")
	if err := s.srv.Shutdown(ctx); err != nil {
		s.log.Errorf("Failed to shutdown server: %v", err)
	}
}
