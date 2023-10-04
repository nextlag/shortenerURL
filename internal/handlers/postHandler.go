package handlers

import (
	"fmt"
	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/storage"
	"io"
	"math/rand"
	"net/http"
)

func PostHandler(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "bad request 400", http.StatusBadRequest)
			return
		}
		shortURL := generateRandomString(8)

		w.WriteHeader(http.StatusCreated)

		_, err = fmt.Fprintf(w, "%s/%s", config.AFlag.URLShort, shortURL)
		if err != nil {
			return
		}
		db.Put(shortURL, string(body))
	}
}

func generateRandomString(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
