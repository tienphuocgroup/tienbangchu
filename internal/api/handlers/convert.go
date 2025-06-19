package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"vietnamese-converter/pkg/converter"
	"vietnamese-converter/pkg/logger"
)

type ConvertResponse struct {
	Number         int64   `json:"number"`
	Vietnamese     string  `json:"vietnamese"`
	ProcessingTimeMs float64 `json:"processing_time_ms"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

type ConvertHandler struct {
	converter converter.NumberConverter
	logger    logger.Logger
}

func (h *ConvertHandler) sendError(w http.ResponseWriter, statusCode int, message, details string) {
	w.WriteHeader(statusCode)
	err := ErrorResponse{
		Error:   message,
		Details: details,
	}
	json.NewEncoder(w).Encode(err)
}

func (h *ConvertHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func (h *ConvertHandler) ConvertNumber(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	var req struct {
		Number   int64  `json:"number"`
		Currency string `json:"currency,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Set default currency if not provided
	if req.Currency == "" {
		req.Currency = "đồng"
	}

	// Validate input
	if req.Number < 0 {
		h.sendError(w, http.StatusBadRequest, "Number must be non-negative", "")
		return
	}

	if req.Number > 999999999999999 {
		h.sendError(w, http.StatusBadRequest, "Number too large", "Maximum supported: 999,999,999,999,999")
		return
	}

	// Convert number
	vietnamese, err := h.converter.ConvertWithCurrency(req.Number, req.Currency)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Conversion failed: %v", err))
		if err.Error() == "number too large (max: 999,999,999,999,999)" || err.Error() == "negative numbers not supported" {
			h.sendError(w, http.StatusBadRequest, "Invalid number", err.Error())
		} else {
			// For other unexpected errors from converter (e.g. potential panics if not caught by middleware)
			h.sendError(w, http.StatusInternalServerError, "Conversion failed unexpectedly", err.Error())
		}
		return
	}

	// Calculate processing time
	processingTime := float64(time.Since(startTime).Nanoseconds()) / 1e6

	// Send response
	response := ConvertResponse{
		Number:          req.Number,
		Vietnamese:      vietnamese,
		ProcessingTimeMs: processingTime,
	}

	h.logger.WithField("number", strconv.FormatInt(req.Number, 10)).
		WithField("processing_time_ms", fmt.Sprintf("%.2f", processingTime)).
		Info("Number converted successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func NewConvertHandler(converter converter.NumberConverter, logger logger.Logger) *ConvertHandler {
	return &ConvertHandler{
		converter: converter,
		logger:    logger,
	}
}

func (h *ConvertHandler) ConvertFromURL(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Get query parameters
	numberStr := r.URL.Query().Get("number")
	if numberStr == "" {
		h.sendError(w, http.StatusBadRequest, "Missing number parameter", "")
		return
	}

	currency := r.URL.Query().Get("currency")
	if currency == "" {
		currency = "đồng" // Default currency
	}

	number, err := strconv.ParseInt(numberStr, 10, 64)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid number format", err.Error())
		return
	}

	// Validate input
	if number < 0 {
		h.sendError(w, http.StatusBadRequest, "Number must be non-negative", "")
		return
	}

	if number > 999999999999999 {
		h.sendError(w, http.StatusBadRequest, "Number too large", "Maximum supported: 999,999,999,999,999")
		return
	}

	// Convert number
	vietnamese, err := h.converter.ConvertWithCurrency(number, currency)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Conversion failed: %v", err))
		if err.Error() == "number too large (max: 999,999,999,999,999)" || err.Error() == "negative numbers not supported" {
			h.sendError(w, http.StatusBadRequest, "Invalid number", err.Error())
		} else {
			// For other unexpected errors from converter (e.g. potential panics if not caught by middleware)
			h.sendError(w, http.StatusInternalServerError, "Conversion failed unexpectedly", err.Error())
		}
		return
	}

	// Calculate processing time
	processingTime := float64(time.Since(startTime).Nanoseconds()) / 1e6

	// Send response
	response := ConvertResponse{
		Number:          number,
		Vietnamese:      vietnamese,
		ProcessingTimeMs: processingTime,
	}

	h.logger.WithField("number", strconv.FormatInt(number, 10)).
		WithField("processing_time_ms", fmt.Sprintf("%.2f", processingTime)).
		Info("Number converted successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}