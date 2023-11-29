package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/service/app"
	"github.com/nextlag/shortenerURL/internal/storage/database/dbstorage"
	"github.com/nextlag/shortenerURL/internal/utils/lg"
)

// Save - обработчик POST-запросов для создания и сохранения URL в storage.
func Save(db app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := lg.New()
		// Считываем тело запроса (оригинальный URL)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "bad request 400", http.StatusBadRequest)
			return
		}

		//  Попытка сохранить short-URL и оригинальный URL в хранилище
		alias, err := db.Put(r.Context(), string(body))
		// обработка конфликта дубликатов
		if errors.Is(err, dbstorage.ErrConflict) {
			log.Error("duplicate url", zap.String("alias", alias), zap.String("url", string(body)))
			http.Error(w, fmt.Sprintf("%s/%s", config.Config.URLShort, alias), http.StatusConflict)
			return
		}

		if err != nil {
			http.Error(w, fmt.Sprintf("failed to add URL: %s", err), http.StatusInternalServerError)
			return
		}

		// Устанавливаем статус HTTP 201 Created
		w.WriteHeader(http.StatusCreated)

		// Отправляем short-URL в теле HTTP-ответа
		_, err = fmt.Fprintf(w, "%s/%s", config.Config.URLShort, alias)
		if err != nil {
			log.Error("error sending short URL response", zap.Error(err))
			return
		}

		// Добавим код для успешного случая
		log.Info("URL added success", zap.String("alias", alias))
	}
}
