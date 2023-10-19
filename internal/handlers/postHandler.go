package handlers

import (
	"fmt"
	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/lib/generateString"
	"github.com/nextlag/shortenerURL/internal/storage"
	"io"
	"net/http"
)

type Request struct {
	Alias string `json:"alias,omitempty"`
}

const aliasLength = 8

// PostHandler - обработчик POST-запросов для создания и сохранения URL в storage.
func PostHandler(urlSaver storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request

		// Считываем тело запроса (оригинальный URL)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "bad request 400", http.StatusBadRequest)
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = generateString.GenerateRandomString(aliasLength)
		}

		// Попытка сохранить short-URL и оригинальный URL в хранилище
		if err := urlSaver.SaveURL(alias, string(body)); err != nil {
			http.Error(w, "internal server error 500", http.StatusInternalServerError)
			return
		}

		// Устанавливаем статус HTTP 201 Created
		w.WriteHeader(http.StatusCreated)

		// Отправляем short-URL в теле HTTP-ответа
		_, err = fmt.Fprintf(w, "%s/%s", config.Args.URLShort, alias)
		if err != nil {
			_ = fmt.Errorf("error sending short URL response: %v", err)
			return
		}
	}
}
