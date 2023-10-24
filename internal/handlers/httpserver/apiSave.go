package httpserver

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/nextlag/shortenerURL/internal/config"
	resp "github.com/nextlag/shortenerURL/internal/lib/api/response"
	"github.com/nextlag/shortenerURL/internal/lib/generatestring"
	"github.com/nextlag/shortenerURL/internal/storage"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	Result string `json:"result"`
}

const aliasLength = 8

func Shorten(log *zap.Logger, db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		// Проверяем, что Content-Type равен "application/json"
		if contentType != "application/json" {
			http.Error(w, "expected Content-Type: application/json", http.StatusBadRequest)
			return
		} else {
			var req Request
			err := render.DecodeJSON(r.Body, &req)
			if errors.Is(err, io.EOF) {
				// Такую ошибку встретим, если получили запрос с пустым телом.
				log.Error("request body is empty")
				render.JSON(w, r, resp.Error("empty request"))
				return
			}
			if err != nil {
				log.Error("failed to decode request body")
				render.JSON(w, r, resp.Error("failed to decode request"))
				return
			}

			log.Info("request body decoded", zap.Any("request", req))
			if err := validator.New().Struct(req); err != nil {
				var validateErr validator.ValidationErrors
				errors.As(err, &validateErr)
				log.Error("invalid request")
				render.JSON(w, r, resp.ValidationError(validateErr))
				return
			}

			alias := req.Alias
			if alias == "" {
				alias = generatestring.NewRandomString(aliasLength)
			}

			id := db.Put(alias, req.URL)
			if id != nil {
				log.Error("failed to add url")
				render.JSON(w, r, resp.Error("failed to add url"))
				return
			}
			responseCreated(w, r, alias)
		}

	}
}

func responseCreated(w http.ResponseWriter, r *http.Request, alias string) {
	response := Response{
		Result: config.Args.URLShort + "/" + alias,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Обработка ошибки кодирования JSON
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
	}
}
