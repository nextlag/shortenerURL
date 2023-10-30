package main_test

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/handlers/httpserver"
	gz "github.com/nextlag/shortenerURL/internal/middleware/gzip"
	"github.com/nextlag/shortenerURL/internal/storage"
)

func TestMain(m *testing.M) {
	if err := config.InitializeArgs(); err != nil {
		log.Fatal(err)
	}

	flag.Parse()
	exitCode := m.Run()
	// Закрываем все оставшиеся тела ответов
	http.DefaultTransport.(*http.Transport).CloseIdleConnections()
	// Завершаем выполнение программы с кодом завершения.
	os.Exit(exitCode)
}
func TestGetHandler(t *testing.T) {
	// Создаем фейковое хранилище
	db := storage.New()
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
			httpserver.GetHandler(db).ServeHTTP(w, req)
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
			db := storage.New()
			// Создаем объект reqBody, который реализует интерфейс io.Reader и будет представлять тело запроса.
			reqBody := strings.NewReader(test.body)
			// Создаем новый POST запрос с текстовым телом и Content-Type: text/plain
			req := httptest.NewRequest("POST", "/", reqBody)
			req.Header.Set("Content-Type", "text/plain")
			// Создаем записывающий ResponseRecorder, который будет использоваться для записи HTTP ответа.
			w := httptest.NewRecorder()
			// Вызываем обработчик для HTTP POST запроса
			httpserver.Save(db).ServeHTTP(w, req)
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

func TestGzipMiddleware(t *testing.T) {
	// Структура с параметрами для тестовых случаев.
	testCases := []struct {
		name              string
		acceptEncoding    string
		expectedBodyEmpty bool
		expectedHeader    string
	}{
		{
			name:           "Accepts_Gzip_Encoding",
			acceptEncoding: "gzip",
			// В этом случае ожидается, что тело ответа будет сжато.
			expectedBodyEmpty: false,
			expectedHeader:    "TestHeader",
		},
		{
			name:           "No_Gzip_Encoding",
			acceptEncoding: "",
			// В этом случае ожидается, что тело ответа будет пустым.
			expectedBodyEmpty: true,
			expectedHeader:    "TestHeader",
		},
	}

	// Создаем тестовый HTTP-обработчик для тестирования middleware.
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Custom-Header", "TestHeader")
		w.WriteHeader(http.StatusOK)
	})

	// Создаем экземпляр middleware с тестовым обработчиком.
	middleware := gz.NewGzip(testHandler)

	// Проходим по всем тестовым случаям.
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Создаем тестовый запрос с указанным заголовком Accept-Encoding.
			req := httptest.NewRequest("GET", "http://example.com", nil)
			req.Header.Set("Accept-Encoding", tc.acceptEncoding)

			// Создаем тестовую запись (Recorder) для записи ответа.
			rr := httptest.NewRecorder()

			// Запускаем middleware.
			middleware.ServeHTTP(rr, req)

			// Проверяем, что middleware корректно обработал запрос.
			if rr.Code != http.StatusOK {
				t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
			}

			// Проверяем, что middleware вернул ожидаемое тело ответа.
			if (len(rr.Body.Bytes()) == 0) != tc.expectedBodyEmpty {
				t.Errorf("Expected empty response body: %v, but got: %v", tc.expectedBodyEmpty, len(rr.Body.Bytes()) == 0)
			}

			// Проверяем, что middleware добавил заголовок X-Custom-Header.
			if rr.Header().Get("X-Custom-Header") != tc.expectedHeader {
				t.Errorf("Expected X-Custom-Header to be '%s', but got %s", tc.expectedHeader, rr.Header().Get("X-Custom-Header"))
			}
		})
	}
}
func TestShorten(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		expectedJSON string
	}{
		{
			name:         "ValidRequest",
			body:         `{"url": "https://example.com", "alias": "example"}`,
			expectedJSON: `{"result":"http://localhost:8080/example"}`,
		},
		{
			name:         "Empty Request Body1",
			body:         `{}`,
			expectedJSON: `{"error":"поле URL обязательно для заполнения"}`,
		},
		{
			name:         "Empty Request Body2",
			body:         `{"url": "example.com"}`,
			expectedJSON: `{"error":"поле URL не является допустимым URL"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Создаем фейковое хранилище
			db := storage.New()
			log := zap.NewNop()
			// Создаем объект reqBody, который реализует интерфейс io.Reader и будет представлять тело запроса.
			reqBody := strings.NewReader(test.body)
			// Создаем новый POST запрос
			req := httptest.NewRequest("POST", "/api/shorten", reqBody)
			req.Header.Set("Content-Type", "application/json")
			// Создаем записывающий ResponseRecorder, который будет использоваться для записи HTTP ответа.
			w := httptest.NewRecorder()
			// Вызываем обработчик для HTTP POST запроса
			httpserver.Shorten(log, db).ServeHTTP(w, req)
			// Получаем результат (HTTP-ответ) после выполнения запроса.
			resp := w.Result()
			defer resp.Body.Close()

			// Чтобы сравнить JSON, сначала декодируем его из ответа сервера.
			var responseJSON map[string]interface{}
			err := json.NewDecoder(resp.Body).Decode(&responseJSON)
			require.NoError(t, err)

			// Затем парсим ожидаемый JSON-ответ для сравнения.
			var expectedJSON map[string]interface{}
			err = json.NewDecoder(strings.NewReader(test.expectedJSON)).Decode(&expectedJSON)
			require.NoError(t, err)

			// Проверка, что полученный JSON совпадает с ожидаемым.
			assert.Equal(t, expectedJSON, responseJSON)
		})
	}
}
