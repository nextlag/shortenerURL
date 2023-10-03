package handlers

import (
	"fmt"
	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/storage"
	"io"
	"net/http"
)

func JSONPostHandler(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "bad request 400", http.StatusBadRequest)
			return
		}

		shortURL := generateRandomString(8)

		// Возвращаем только URL в виде строки без JSON-обертки
		urlResponse := *config.URLShort + "/" + shortURL

		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(urlResponse)))
		w.WriteHeader(http.StatusCreated)

		_, err = fmt.Fprint(w, urlResponse)
		if err != nil {
			return
		}

		db.Put(shortURL, string(body))
	}
}
