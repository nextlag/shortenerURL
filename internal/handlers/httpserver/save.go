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

// PostHandler - обработчик POST-запросов для создания и сохранения Url в storage.
func Save(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Считываем тело запроса (оригинальный Url)
		url, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "bad request 400", http.StatusBadRequest)
			return
		}
		// Генерируем случайную строку
		alias := generatestring.NewRandomString(8)
		// Вызываем метод Save для чтения мапы и записи их в файл
		err = db.Save(config.Args.FileStorage, alias, string(url))
		if err != nil {
			return
		}

		// Попытка сохранить short-Url и оригинальный Url в хранилище
		if err := db.Put(alias, string(url)); err != nil {
			http.Error(w, "internal server error 500", http.StatusInternalServerError)
			return
		}

		// Устанавливаем статус HTTP 201 Created
		w.WriteHeader(http.StatusCreated)

		// Отправляем short-Url в теле HTTP-ответа
		_, err = fmt.Fprintf(w, "%s/%s", config.Args.URLShort, alias)
		if err != nil {
			log.Printf("Error sending short Url response: %v", err)
			return
		}
	}
}
