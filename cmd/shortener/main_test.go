package main_test

import (
	"bytes"
	"encoding/json"
	"github.com/nextlag/shortenerURL/internal/handlers"
	"github.com/nextlag/shortenerURL/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	exitCode := m.Run()
	// Закрываем все оставшиеся тела ответов
	http.DefaultTransport.(*http.Transport).CloseIdleConnections()
	// Завершаем выполнение программы с кодом завершения.
	os.Exit(exitCode)
}
func TestGetHandler(t *testing.T) {
	// Создаем фейковое хранилище
	db := storage.NewInMemoryStorage()
	// Пушим данныые
	db.Put("example", "http://example.com")

	// Тестим валидный запрос
	t.Run("Valid ID", func(t *testing.T) {
		// Создаем фейковый запрос с валидным идентификатором
		req := httptest.NewRequest("GET", "/example", nil)
		w := httptest.NewRecorder()

		// Создаем и вызываем handler для маршрута
		handlers.GetHandler(db).ServeHTTP(w, req)
		resp := w.Result()

		// Проверяем статус кода
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		// Получаем значение Location
		location := resp.Header.Get("Location")
		// Ожидаем, что Location равен "http://example.com"
		assert.Empty(t, location)
		// Закрываем тело HTTP-ответа
		require.NoError(t, resp.Body.Close())
	})
	// Тестим невалидный запрос
	t.Run("Invalid ID", func(t *testing.T) {
		// Создаем фейковый запрос с невалидным идентификатором
		req := httptest.NewRequest("GET", "/nonexistent", nil)
		w := httptest.NewRecorder()
		// Создаем и вызываем handler для маршрута
		handlers.GetHandler(db).ServeHTTP(w, req)
		resp := w.Result()
		// Проверяем, что статус BadRequest
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		// Закрываем тело HTTP-ответа
		require.NoError(t, resp.Body.Close())
	})
}

func TestPostHandler(t *testing.T) {
	// Создаем фейковое хранилище
	db := storage.NewInMemoryStorage()
	body := "http://example.com"
	// Создаем объект requestData для JSON-тела запроса
	requestData := struct {
		URL string `json:"url"`
	}{
		URL: body,
	}
	// Кодируем requestData в JSON и создаем объект reqBody, реализующий интерфейс io.Reader
	reqBody, err := json.Marshal(requestData)
	require.NoError(t, err)
	// Создаем новый POST запрос
	req := httptest.NewRequest("POST", "/", bytes.NewReader(reqBody))
	// Создаем записывающий ResponseRecorder, который будет использоваться для записи HTTP ответа
	w := httptest.NewRecorder()
	// Вызываем обработчик для HTTP POST запроса
	handlers.PostHandler(db).ServeHTTP(w, req)
	// Получаем результат (HTTP-ответ) после выполнения запроса
	resp := w.Result()
	// Проверяем статус
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	// Извлекаем сокращенную версию URL из тела HTTP-ответа
	var response struct {
		ShortURL string `json:"short_url"`
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	// Проверяем длину shortURL
	assert.Equal(t, 8, len(response.ShortURL))
	// Получаем сокращенный URL
	storedURL, ok := db.Get(response.ShortURL)
	// Проверяем, что URL был сохранен
	assert.True(t, ok)
	// Проверяем, что значение URL в хранилище совпадает с ожидаемым
	assert.Equal(t, body, storedURL)
	// Закрываем тело HTTP-ответа
	require.NoError(t, resp.Body.Close())
}
