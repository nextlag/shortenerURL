package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/service/app"
	"github.com/nextlag/shortenerURL/internal/service/auth"
	"github.com/nextlag/shortenerURL/internal/storage"
	"github.com/nextlag/shortenerURL/internal/storage/dbstorage"
)

type ShortenHandler struct {
	db  app.Storage
	log *zap.Logger
	cfg config.Args
}

func NewShortenHandlers(db app.Storage, log *zap.Logger, cfg config.Args) *ShortenHandler {
	return &ShortenHandler{
		db:  db,
		log: log,
		cfg: cfg,
	}
}

// Shorten - это обработчик HTTP-запросов для сокращения URL.
func (h *ShortenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req storage.FileStorage
	// декодирование JSON-запроса из тела HTTP-запроса в структуру Data.
	err := render.DecodeJSON(r.Body, &req)

	// Обработка случая, когда тело запроса пустое.
	if errors.Is(err, io.EOF) {
		h.log.Error("request body is empty", zap.Error(err))
		render.JSON(w, r, Error("empty request"))
		return
	}

	// Обработка ошибок декодирования тела запроса.
	if err != nil {
		h.log.Error("failed to decode request body", zap.Error(err))
		render.JSON(w, r, Error("failed to decode request"))
		return
	}

	// Валидация входных данных с использованием библиотеки валидатора.
	if err := validator.New().Struct(req); err != nil {
		var validateErr validator.ValidationErrors
		errors.As(err, &validateErr)
		h.log.Error("invalid request")
		render.JSON(w, r, ValidationError(validateErr))
		return
	}

	uuid, err := auth.CheckCookie(w, r, h.log)
	if err != nil {
		h.log.Error("Error getting cookie: ", zap.Error(err))
		return
	}

	// Добавление URL в хранилище и получение идентификатора (id).
	alias, err := h.db.Put(r.Context(), req.URL, uuid)
	if errors.Is(err, dbstorage.ErrConflict) {
		// ошибка для случая конфликта оригинальных url
		h.log.Error("trying to add a duplicate URL", zap.Error(err))
		ResponseConflict(w, alias)
		return
	}

	// Обработка ошибки при добавлении URL в хранилище.
	if err != nil {
		er := fmt.Sprintf("failed to add URL: %s", err)
		render.JSON(w, r, Error(er))
		return
	}
	// Отправка ответа клиенту с сокращенной ссылкой.
	ResponseCreated(w, alias)
}
