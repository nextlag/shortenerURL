package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/render"
	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/service/app"
	"github.com/nextlag/shortenerURL/internal/service/auth"
)

// BatchHandler представляет хендлер для сокращения нескольких URL.
type BatchHandler struct {
	db  app.Storage
	log *zap.Logger
	cfg config.Args
}

// NewBatchHandler создает новый экземпляр BatchHandler.
func NewBatchHandler(db app.Storage, log *zap.Logger, cfg config.Args) *BatchHandler {
	return &BatchHandler{
		db:  db,
		log: log,
		cfg: cfg,
	}
}

type BatchShortenRequest []struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// BatchShortenResponse представляет структуру ответа для сокращения нескольких URL.
type BatchShortenResponse []struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// ServeHTTP обрабатывает HTTP-запрос для сокращения нескольких URL.
func (h *BatchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req BatchShortenRequest
	var resp BatchShortenResponse

	err := render.DecodeJSON(r.Body, &req)

	if errors.Is(err, io.EOF) {
		h.log.Error("request body is empty")
		render.JSON(w, r, Error("empty request"))
		return
	}

	if err != nil {
		h.log.Error("failed to decode request body", zap.Error(err))
		render.JSON(w, r, Error("failed to decode request"))
		return
	}
	uid := auth.CheckCookie(w, r, h.log)

	for _, url := range req {
		alias, err := h.db.Put(r.Context(), url.OriginalURL, uid)
		if err != nil {
			h.log.Error("URL", zap.Error(err))
			return
		}

		resp = append(resp, struct {
			CorrelationID string `json:"correlation_id"`
			ShortURL      string `json:"short_url"`
		}{
			CorrelationID: url.CorrelationID,
			ShortURL:      fmt.Sprintf("%s/%s", h.cfg.BaseURL, alias),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to encode JSON response", http.StatusInternalServerError)
	}
}
