package httpserver

import (
	"encoding/json"
	"errors"
	"fmt"
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

// Request представляет структуру входящего JSON-запроса.
type Request struct {
	URL   string `json:"url" validate:"required,url"` // URL, который нужно сократить, должен быть валидным URL.
	Alias string `json:"alias,omitempty"`             // Alias, Пользовательский псевдоним для короткой ссылки (необязательный).
}

// Response представляет структуру ответа, отправляемого клиенту.
type Response struct {
	Result string `json:"result"` // Результат - сокращенная ссылка.
}

// aliasLength - длина по умолчанию для генерируемых псевдонимов.
const aliasLength = 8

// Shorten - это обработчик HTTP-запросов для сокращения URL.
func Shorten(log *zap.Logger, db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
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

		log.Info("request body decoded", zap.Any("request", req))

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
			log.Error("failed to add URL")
			render.JSON(w, r, resp.Error("failed to add URL"))
			return
			//} else {
			//	err = db.Save(config.Args.FileStorage, alias, req.URL)
			//	if err != nil {
			//		return
			//	}
			//}
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
