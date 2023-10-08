package handlers

import (
	"fmt"
	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/storage"
	"github.com/nextlag/shortenerURL/internal/util"
	"io"
	"log"
	"net/http"
)

// PostHandler - обработчик POST-запросов для создания и сохранения URL в storage.
func PostHandler(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Считываем тело запроса (оригинальный URL)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "bad request 400", http.StatusBadRequest)
			return
		}
		// Генерируем случайную строку
		shortURL := util.GenerateRandomString(8)

		// Попытка сохранить short-URL и оригинальный URL в хранилище
		if err := db.Put(shortURL, string(body)); err != nil {
			http.Error(w, "internal server error 500", http.StatusInternalServerError)
			return
		}

		// Устанавливаем статус HTTP 201 Created
		w.WriteHeader(http.StatusCreated)

		// Отправляем short-URL в теле HTTP-ответа
		_, err = fmt.Fprintf(w, "%s/%s", config.Args.BaseURL, shortURL)
		if err != nil {
			log.Printf("Error sending short URL response: %v", err)
			return
		}
	}
}
