package handlers

import (
	"encoding/json"
	"github.com/nextlag/shortenerURL/internal/storage"
	"math/rand"
	"net/http"
)

func PostHandler(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Чтение JSON-тела запроса
		var requestData struct {
			URL string `json:"url"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// Генерация короткого URL
		shortURL := generateRandomString(8)

		// Сохранение соответствия в хранилище
		db.Put(shortURL, requestData.URL)

		// Отправка короткого URL в ответе
		response := struct {
			ShortURL string `json:"short_url"`
		}{
			ShortURL: shortURL,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
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
