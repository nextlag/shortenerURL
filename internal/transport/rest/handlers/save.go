package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/service/app"
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

		// Попытка сохранить short-URL и оригинальный URL в хранилище
		alias, err := db.Put(string(body))
		if err != nil {
			er := fmt.Sprintf("failed to add URL: %s", err)
			http.Error(w, er, http.StatusInternalServerError)
			return
		}

		// Устанавливаем статус HTTP 201 Created
		w.WriteHeader(http.StatusCreated)

		// Отправляем short-URL в теле HTTP-ответа
		_, err = fmt.Fprintf(w, "%s/%s", config.Config.URLShort, alias)
		if err != nil {
			log.Printf("Error sending short URL response: %v", err)
			return
		}

		// Добавим код для успешного случая
		log.Printf("URL added successfully. Alias: %s", alias)
	}
}
