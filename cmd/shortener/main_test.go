package main_test

import (
	"flag"
	"github.com/nextlag/shortenerURL/internal/handlers/http-server"
	"github.com/nextlag/shortenerURL/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	flag.Parse()
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
	err := db.Put("example", "http://example.com")
	if err != nil {
		return
	}

	tests := []struct {
		Name             string
		RequestPath      string
		ExpectedStatus   int
		ExpectedLocation string
	}{
		{
			Name:             "Valid ID",
			RequestPath:      "/example",
			ExpectedStatus:   http.StatusBadRequest,
			ExpectedLocation: "http://example.com",
		},
		{
			Name:             "Invalid ID",
			RequestPath:      "/nonexistent",
			ExpectedStatus:   http.StatusBadRequest,
			ExpectedLocation: "",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// Создаем фейковый запрос
			req := httptest.NewRequest("GET", test.RequestPath, nil)
			w := httptest.NewRecorder()

			// Создаем и вызываем handler для маршрута
			http_server.GetHandler(db).ServeHTTP(w, req)
			resp := w.Result()

			// Проверяем статус кода
			assert.Equal(t, test.ExpectedStatus, resp.StatusCode)
			// Получаем значение Location
			location := resp.Header.Get("Location")
			assert.Empty(t, location)

			// Закрываем тело HTTP-ответа
			require.NoError(t, resp.Body.Close())
		})
	}
}

func TestTextPostHandler(t *testing.T) {
	tests := []struct {
		name                   string
		body                   string
		expectedURL            string
		expectedShortURLLength int
	}{
		{
			name:                   "ValidRequest",
			body:                   "http://example.com",
			expectedURL:            "http://localhost:8080/",
			expectedShortURLLength: 8,
		},
		{
			name:                   "LongURL",
			body:                   "https://www.thisisaverylongurlthatexceedsthecharacterlimit.com",
			expectedURL:            "http://localhost:8080/",
			expectedShortURLLength: 8,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Создаем фейковое хранилище
			db := storage.NewInMemoryStorage()
			// Создаем объект reqBody, который реализует интерфейс io.Reader и будет представлять тело запроса.
			reqBody := strings.NewReader(test.body)
			// Создаем новый POST запрос с текстовым телом и Content-Type: text/plain
			req := httptest.NewRequest("POST", "/", reqBody)
			req.Header.Set("Content-Type", "text/plain")
			// Создаем записывающий ResponseRecorder, который будет использоваться для записи HTTP ответа.
			w := httptest.NewRecorder()
			// Вызываем обработчик для HTTP POST запроса
			http_server.Save(db).ServeHTTP(w, req)
			// Получаем результат (HTTP-ответ) после выполнения запроса.
			resp := w.Result()
			defer resp.Body.Close()

			// Проверяем статус
			assert.Equal(t, http.StatusCreated, resp.StatusCode)

			// Извлекаем сокращенную версию URL из тела HTTP-ответа, удаляя из неё префикс ожидаемого URL.
			shortURL := strings.TrimPrefix(w.Body.String(), test.expectedURL)
			// Проверяем длину shortURL
			assert.Equal(t, test.expectedShortURLLength, len(shortURL))

			// Получаем сокращенный URL из хранилища
			storedURL, err := db.Get(shortURL)
			// Проверяем, что нет ошибки при получении URL из хранилища
			require.NoError(t, err)
			// Проверяем, что значение URL в хранилище не пустое
			assert.NotEmpty(t, storedURL)
		})
	}
}
