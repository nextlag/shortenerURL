package httpserver

import (
	"fmt"
	"github.com/go-chi/render"
	"github.com/nextlag/shortenerURL/internal/config"
	resp "github.com/nextlag/shortenerURL/internal/lib/api/response"
	"github.com/nextlag/shortenerURL/internal/lib/generatestring"
	"github.com/nextlag/shortenerURL/internal/storage"
	"io"
	"log"
	"net/http"
)

// Save - обработчик POST-запросов для создания и сохранения URL в storage.
func Save(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Считываем тело запроса (оригинальный URL)
		url, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "bad request 400", http.StatusBadRequest)
			return
		}
		// Генерируем случайную строку
		alias := generatestring.NewRandomString(8)

		// Попытка сохранить short-URL и оригинальный URL в хранилище
		err = db.Put(alias, string(url))
		if err != nil {
			er := fmt.Sprintf("failed to add URL: %s", err)
			render.JSON(w, r, resp.Error(er))
			return
		} else {
			err := db.Save(config.Args.FileStorage, alias, string(url))
			if err != nil {
				er := fmt.Sprintf("failed to add URL: %s", err)
				render.JSON(w, r, resp.Error(er))
			}
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
