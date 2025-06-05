package middleware

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"vietnamese-converter/pkg/logger"

	"github.com/google/uuid"
	"golang.org/x/time/rate"
)

func RequestLogger(logger logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			
			next.ServeHTTP(wrapped, r)
			
			duration := time.Since(start)
			
			logger.WithField("method", r.Method).
				WithField("path", r.URL.Path).
				WithField("status", fmt.Sprintf("%d", wrapped.statusCode)).
				WithField("duration_ms", fmt.Sprintf("%.2f", float64(duration.Nanoseconds())/1e6)).
				WithField("remote_addr", r.RemoteAddr).
				Info("HTTP request processed")
		})
	}
}

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		ctx := context.WithValue(r.Context(), "request_id", requestID)
		w.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Recoverer(logger logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error(fmt.Sprintf("Panic recovered: %v\n%s", err, debug.Stack()))
					
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error":"Internal Server Error","code":500}`))
				}
			}()
			
			next.ServeHTTP(w, r)
		})
	}
}

func RateLimiter(requestsPerSecond int) func(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Limit(requestsPerSecond), requestsPerSecond)
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"error":"Rate limit exceeded","code":429}`))
				return
			}
			
			next.ServeHTTP(w, r)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
