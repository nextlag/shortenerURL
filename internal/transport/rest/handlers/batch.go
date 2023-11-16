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
	"github.com/nextlag/shortenerURL/internal/utils/generatestring"
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

// type BatchShortenRequest []struct {
// 	ID  string `json:"correlation_id"`
// 	URL string `json:"original_url"`
// }
//
// // BatchShortenResponse представляет структуру ответа для сокращения нескольких URL.
// type BatchShortenResponse []struct {
// 	ID    string `json:"correlation_id"`
// 	Alias string `json:"short_url"`
// }
//
// // ServeHTTP обрабатывает HTTP-запрос для сокращения нескольких URL.
// func (h *BatchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	var req BatchShortenRequest
// 	err := render.DecodeJSON(r.Body, &req)
//
// 	if errors.Is(err, io.EOF) {
// 		h.log.Error("request body is empty")
// 		render.JSON(w, r, resp.Error("empty request"))
// 		return
// 	}
//
// 	if err != nil {
// 		h.log.Error("failed to decode request body", zap.Error(err))
// 		render.JSON(w, r, resp.Error("failed to decode request"))
// 		return
// 	}
//
// 	var response BatchShortenResponse
//
// 	for _, url := range req {
// 		alias := generatestring.NewRandomString(aliasLength)
//
// 		if err := h.db.Put(alias, url.URL); err != nil {
// 			er := "failed to add URL: " + err.Error()
// 			render.JSON(w, r, resp.Error(er))
// 			return
// 		}
//
// 		response = append(response, struct {
// 			ID    string `json:"correlation_id"`
// 			Alias string `json:"short_url"`
// 		}{
// 			ID:    url.ID,
// 			Alias: fmt.Sprintf("%s/%s", config.Config.URLShort, alias),
// 		})
// 	}
//
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
//
// 	if err := json.NewEncoder(w).Encode(response); err != nil {
// 		http.Error(w, "failed to encode JSON response", http.StatusInternalServerError)
// 	}
// }

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
	h.log.Info("Request Body", zap.Any("body", r.Body))
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
		alias := generatestring.NewRandomString(aliasLength)

		if err := h.db.Put(alias, url.OriginalURL); err != nil {
			er := "failed to add URL: " + err.Error()
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
