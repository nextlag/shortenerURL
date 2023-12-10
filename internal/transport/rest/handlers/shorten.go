package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/service/app"
	"github.com/nextlag/shortenerURL/internal/service/auth"
	"github.com/nextlag/shortenerURL/internal/storage/dbstorage"
	"github.com/nextlag/shortenerURL/internal/usecase"
)

// Shorten - это обработчик HTTP-запросов для сокращения URL.
func Shorten(log *zap.Logger, db app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req usecase.CustomRequest
		// декодирование JSON-запроса из тела HTTP-запроса в структуру Request.
		err := render.DecodeJSON(r.Body, &req)

		// Обработка случая, когда тело запроса пустое.
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty", zap.Error(err))
			render.JSON(w, r, Error("empty request"))
			return
		}

		// Обработка ошибок декодирования тела запроса.
		if err != nil {
			log.Error("failed to decode request body", zap.Error(err))
			render.JSON(w, r, Error("failed to decode request"))
			return
		}

		// Валидация входных данных с использованием библиотеки валидатора.
		if err := validator.New().Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)
			log.Error("invalid request")
			render.JSON(w, r, ValidationError(validateErr))
			return
		}

		uid := auth.CheckCookie(w, r, log)

		// Добавление URL в хранилище и получение идентификатора (id).
		alias, err := db.Put(r.Context(), req.GetEntityRequest().URL, uid)
		if errors.Is(err, dbstorage.ErrConflict) {
			// ошибка для случая конфликта оригинальных url
			log.Error("Извините, такой url уже занят")
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
}
