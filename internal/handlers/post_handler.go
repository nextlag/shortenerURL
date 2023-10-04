package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/storage"
	"io"
	"math/rand"
	"net/http"
)

func PostHandler(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		if contentType != "text/plain" && contentType != "application/json" {
			http.Error(w, "unsupported media type 400", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "bad request 400", http.StatusBadRequest)
			return
		}

		shortURL := generateRandomString(8)

		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", "30")
		w.WriteHeader(http.StatusCreated)

		if contentType == "application/json" {
			// Предполагается, что JSON имеет структуру {"data": "значение"}
			var jsonData struct {
				Data string `json:"data"`
			}
			if err := json.Unmarshal(body, &jsonData); err != nil {
				http.Error(w, "bad request 400", http.StatusBadRequest)
				return
			}
			body = []byte(jsonData.Data)
		}
		_, err = fmt.Fprintf(w, "%s/%s", *config.URLShort, shortURL)
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
