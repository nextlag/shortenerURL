package httpserver

import (
	"fmt"
	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/lib/generatestring"
	"github.com/nextlag/shortenerURL/internal/storage"
	"io"
	"log"
	"net/http"
)

func Save(db storage.Storage) http.HandlerFunc {
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
			http.Error(w, "internal server error 500", http.StatusInternalServerError)
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
