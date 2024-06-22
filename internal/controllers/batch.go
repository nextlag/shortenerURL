// Package controllers provides the handlers for managing URL shortening operations.
package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/render"
	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/usecase/auth"
)

// BatchShortenRequest represents a request structure for shortening multiple URLs.
type BatchShortenRequest []struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// BatchShortenResponse represents a response structure for shortening multiple URLs.
type BatchShortenResponse []struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// Batch handles the HTTP request for shortening multiple URLs.
// It decodes the JSON request, processes each URL, and returns the shortened URLs in the response.
// If an error occurs during processing, it logs the error and responds with an appropriate message.
func (c *Controller) Batch(w http.ResponseWriter, r *http.Request) {
	var req BatchShortenRequest
	var resp BatchShortenResponse

	err := render.DecodeJSON(r.Body, &req)

	if errors.Is(err, io.EOF) {
		c.log.Error("request body is empty")
		render.JSON(w, r, Error("empty request"))
		return
	}

	if err != nil {
		c.log.Error("failed to decode request body", zap.Error(err))
		render.JSON(w, r, Error("failed to decode request"))
		return
	}
	uuid, err := auth.CheckCookie(w, r, c.log)
	if err != nil {
		c.log.Error("Error getting cookie: ", zap.Error(err))
		return
	}

	for _, url := range req {
		alias, err := c.uc.DoPut(r.Context(), url.OriginalURL, uuid)
		if err != nil {
			c.log.Error("URL", zap.Error(err))
			return
		}

		resp = append(resp, struct {
			CorrelationID string `json:"correlation_id"`
			ShortURL      string `json:"short_url"`
		}{
			CorrelationID: url.CorrelationID,
			ShortURL:      fmt.Sprintf("%s/%s", c.cfg.BaseURL, alias),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to encode JSON response", http.StatusInternalServerError)
	}
}
