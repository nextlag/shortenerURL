package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/service/app"
	"github.com/nextlag/shortenerURL/internal/utils/generatestring"
)

// PostHandler - обработчик POST-запросов для создания и сохранения URL в storage.
func Save(db app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Считываем тело запроса (оригинальный URL)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "bad request 400", http.StatusBadRequest)
			return
		}
		// Генерируем случайную строку
		alias := generatestring.NewRandomString(aliasLength)

		// Попытка сохранить short-URL и оригинальный URL в хранилище
		if err := db.Put(alias, string(body)); err != nil {
			er := fmt.Sprintf("failed to add URL: %s", err)
			http.Error(w, er, http.StatusInternalServerError)
			return
		}

		// Устанавливаем статус HTTP 201 Created
		w.WriteHeader(http.StatusCreated)

		// Отправляем short-URL в теле HTTP-ответа
		_, err = fmt.Fprintf(w, "%s/%s", config.Args.URLShort, alias)
		if err != nil {
			log.Printf("Error sending short URL response: %v", err)
			return
		}
	}
}
