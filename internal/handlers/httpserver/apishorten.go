package httpserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/app"
	"github.com/nextlag/shortenerURL/internal/config"
	resp "github.com/nextlag/shortenerURL/internal/lib/api/response"
	"github.com/nextlag/shortenerURL/internal/lib/filestorage"
	"github.com/nextlag/shortenerURL/internal/lib/generatestring"
)

// Response представляет структуру ответа, отправляемого клиенту.
type Response struct {
	Result string `json:"result"` // Результат - сокращенная ссылка.
}

// aliasLength - длина по умолчанию для генерируемых алиасов.
const aliasLength = 8

// Shorten - это обработчик HTTP-запросов для сокращения URL.
func Shorten(log *zap.Logger, db app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req filestorage.Request
		// декодирование JSON-запроса из тела HTTP-запроса в структуру Request.
		err := render.DecodeJSON(r.Body, &req)

		// Обработка случая, когда тело запроса пустое.
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			render.JSON(w, r, resp.Error("empty request"))
			return
		}

		// Обработка ошибок декодирования тела запроса.
		if err != nil {
			log.Error("failed to decode request body")
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		// Валидация входных данных с использованием библиотеки валидатора.
		if err := validator.New().Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)
			log.Error("invalid request")
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		// Генерация алиаса, если пользовательский алиас не указан.
		alias := req.Alias
		if alias == "" {
			alias = generatestring.NewRandomString(aliasLength)
		}
		// Добавление URL в хранилище и получение идентификатора (id).
		err = db.Put(alias, req.URL)
		// Обработка ошибки при добавлении URL в хранилище.
		if err != nil {
			er := fmt.Sprintf("failed to add URL: %s", err)
			render.JSON(w, r, resp.Error(er))
			return
		}
		// Отправка ответа клиенту с сокращенной ссылкой.
		responseCreated(w, alias)
	}
}

// responseCreated отправляет успешный ответ с сокращенной ссылкой в JSON-формате.
func responseCreated(w http.ResponseWriter, alias string) {
	response := Response{
		Result: fmt.Sprintf("%s/%s", config.Args.URLShort, alias),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Обработка ошибки кодирования JSON.
		http.Error(w, "failed to encode JSON response", http.StatusInternalServerError)
	}
}
