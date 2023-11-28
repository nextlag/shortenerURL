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
	resp "github.com/nextlag/shortenerURL/internal/transport/rest/response"
)

// BatchHandler представляет хендлер для сокращения нескольких URL.
type BatchHandler struct {
	log *zap.Logger
	db  app.Storage
}

// NewBatchHandler создает новый экземпляр BatchHandler.
func NewBatchHandler(log *zap.Logger, db app.Storage) *BatchHandler {
	return &BatchHandler{
		log: log,
		db:  db,
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
	err := render.DecodeJSON(r.Body, &req)

	if errors.Is(err, io.EOF) {
		h.log.Error("request body is empty")
		render.JSON(w, r, resp.Error("empty request"))
		return
	}

	if err != nil {
		h.log.Error("failed to decode request body", zap.Error(err))
		render.JSON(w, r, resp.Error("failed to decode request"))
		return
	}

	var response BatchShortenResponse

	for _, url := range req {

		alias, err := h.db.Put(url.OriginalURL)
		if err != nil {
			er := "URL" + err.Error()
			render.JSON(w, r, resp.Error(er))
			return
		}

		response = append(response, struct {
			CorrelationID string `json:"correlation_id"`
			ShortURL      string `json:"short_url"`
		}{
			CorrelationID: url.CorrelationID,
			ShortURL:      fmt.Sprintf("%s/%s", config.Config.URLShort, alias),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode JSON response", http.StatusInternalServerError)
	}
}
