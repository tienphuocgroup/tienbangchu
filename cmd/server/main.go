package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"vietnamese-converter/internal/api/handlers"
	"vietnamese-converter/internal/api/middleware"
	"vietnamese-converter/internal/api/routes"
	"vietnamese-converter/internal/config"
	"vietnamese-converter/pkg/converter"
	"vietnamese-converter/pkg/logger"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.Load()
	logger := logger.New(cfg.Log.Level)
	logger.Info("Starting Vietnamese Number Converter Service")

	vietnameseConverter := converter.NewVietnameseConverter()
	convertHandler := handlers.NewConvertHandler(vietnameseConverter, logger)
	router := setupRouter(convertHandler, logger)
	
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		logger.Info(fmt.Sprintf("Server starting on port %d", cfg.Server.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(fmt.Sprintf("Server failed to start: %v", err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal(fmt.Sprintf("Server forced to shutdown: %v", err))
	}

	logger.Info("Server shutdown complete")
}

func setupRouter(convertHandler *handlers.ConvertHandler, logger logger.Logger) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestLogger(logger))
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer(logger))
	r.Use(middleware.RateLimiter(10000))
	routes.SetupConvertRoutes(r, convertHandler)
	return r
}
