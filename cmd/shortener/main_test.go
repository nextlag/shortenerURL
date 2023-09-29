package main_test

import (
	"github.com/nextlag/shortenerURL/internal/handlers"
	"github.com/nextlag/shortenerURL/internal/storage"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetHandler(t *testing.T) {
	// Создаем фейковое хранилище
	db := storage.NewInMemoryStorage()
	db.Put("example", "http://example.com")

	t.Run("Valid ID", func(t *testing.T) {
		// Создаем фейковый запрос с валидным идентификатором
		req := httptest.NewRequest("GET", "/example", nil)
		// Создаем фейковый ResponseWriter
		w := httptest.NewRecorder()

		// Создаем обработчик для маршрута и вызываем его
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlers.GetHandler(db, w, r)
		})
		handler.ServeHTTP(w, req)

		// Получаем HTTP ответ
		resp := w.Result()

		// Проверяем статус код ответа, он должен быть 307 (Temporary Redirect)
		if resp.StatusCode != http.StatusTemporaryRedirect {
			t.Errorf("Expected status code %d, but got %d", http.StatusTemporaryRedirect, resp.StatusCode)
		}

		// Проверяем заголовок Location
		location := resp.Header.Get("Location")
		if location != "http://example.com" {
			t.Errorf("Expected Location header to be 'http://example.com', but got '%s'", location)
		}
		// Закрываем тело HTTP-ответов
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Print(err)
			}
		}(resp.Body)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		// Создаем фейковый запрос с невалидным идентификатором
		req := httptest.NewRequest("GET", "/nonexistent", nil)
		// Создаем фейковый ResponseWriter
		w := httptest.NewRecorder()

		// Создаем обработчик для маршрута и вызываем его
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlers.GetHandler(db, w, r)
		})
		handler.ServeHTTP(w, req)

		// Получаем HTTP ответ
		resp := w.Result()

		// Проверяем статус код ответа, он должен быть 400 (Bad Request)
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, resp.StatusCode)
		}
	})
}

func TestPostHandler(t *testing.T) {
	// Создаем фейковое хранилище
	db := storage.NewInMemoryStorage()

	// Создаем фейковый запрос с телом
	body := "http://example.com"
	reqBody := strings.NewReader(body)
	req := httptest.NewRequest("POST", "/", reqBody)
	w := httptest.NewRecorder()

	// Вызываем обработчик
	handlers.PostHandler(db, w, req)

	// Проверяем статус код ответа, он должен быть 201 (Created)
	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code %d, but got %d", http.StatusCreated, resp.StatusCode)
	}

	// Проверяем, что тело ответа содержит сгенерированный URL
	expectedURL := "http://localhost:8080/"
	shortURL := strings.TrimPrefix(w.Body.String(), expectedURL)
	if len(shortURL) != 8 {
		t.Errorf("Expected short URL to be of length 8, but got %d", len(shortURL))
	}

	// Проверяем, что данные были добавлены в хранилище
	storedURL, ok := db.Get(shortURL)
	if !ok {
		t.Error("Expected URL to be stored in the database, but it's not")
	}
	if storedURL != body {
		t.Errorf("Expected request body to be saved in the database, but it's not. Got: %s, Expected: %s", storedURL, body)
	}
	// Закрываем тело HTTP-ответов
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Print(err)
		}
	}(resp.Body)
}
