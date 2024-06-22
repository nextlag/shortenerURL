package logger_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/middleware/logger"
)

func TestLoggerMiddleware(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		path     string
		expected string
	}{
		{
			name:     "GET request",
			method:   http.MethodGet,
			path:     "/test",
			expected: "GET /test",
		},
		{
			name:     "POST request",
			method:   http.MethodPost,
			path:     "/data",
			expected: "POST /data",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Создаем тестовый логгер
			testLogger := zap.NewExample()

			// Настраиваем middleware логгера
			cfg := config.HTTPServer{FileStorage: "/path/to/storage"}
			mw := logger.New(testLogger, cfg)

			// Создаем тестовый сервер с middleware логгера
			handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Проверяем, что данные запроса совпадают с ожидаемыми
				assert.Equal(t, tc.method, r.Method)
				assert.Equal(t, tc.path, r.URL.Path)

				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			}))

			req := httptest.NewRequest(tc.method, tc.path, nil)
			req.Header.Set("Content-Type", "application/json")
			req.RemoteAddr = "127.0.0.1:12345"
			req.Header.Set("User-Agent", "Test User Agent")

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			// Проверяем, что запрос был обработан успешно
			assert.Equal(t, http.StatusOK, rec.Code)
		})
	}
}

func BenchmarkNewLoggerMiddleware(b *testing.B) {
	// Создаем тестовый логгер
	testLogger := zap.NewNop()

	// Создаем конфигурацию HTTP сервера
	cfg := config.HTTPServer{FileStorage: "/path/to/storage"}

	// Настраиваем запрос и обработчик для бенчмарка
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	// Запускаем бенчмарк
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Создаем middleware
		mw := logger.New(testLogger, cfg)

		// Запускаем middleware
		rr := &fakeResponseWriter{}
		mw(handler).ServeHTTP(rr, req)
	}
}

// fakeResponseWriter реализует http.ResponseWriter для тестов
type fakeResponseWriter struct {
	header http.Header
	code   int
}

func (rw *fakeResponseWriter) Header() http.Header {
	if rw.header == nil {
		rw.header = make(http.Header)
	}
	return rw.header
}

func (rw *fakeResponseWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

func (rw *fakeResponseWriter) WriteHeader(statusCode int) {
	rw.code = statusCode
}

// Example of using the logging middleware in an HTTP server.
func ExampleNew() {
	// Initialize the logger
	log := logger.SetupLogger()

	// Define a simple handler that returns a plain text response.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, logged world!"))
	})

	// Create a new router and add the logging middleware.
	r := chi.NewRouter()
	r.Use(logger.New(log, config.HTTPServer{
		FileStorage: "/tmp/data",
	}))

	// Add the handler to the router.
	r.Get("/", handler)

	// Create a new HTTP request.
	req := httptest.NewRequest("GET", "http://example.com", nil)
	w := httptest.NewRecorder()

	// Serve the HTTP request using the router with logging middleware.
	r.ServeHTTP(w, req)

	// Get the recorded response.
	res := w.Result()
	defer res.Body.Close()

	// Print the response status and body.
	fmt.Println("Status:", res.StatusCode)
	body, _ := io.ReadAll(res.Body)
	fmt.Println("Body:", string(body))

	// Output:
	// Status: 200
	// Body: Hello, logged world!
}

func ExampleSetupLogger() {
	// Initialize the logger
	log := logger.SetupLogger()

	// Log a simple info message
	log.Info("example log message", zap.String("exampleKey", "exampleValue"))

	// Output:
	// {"level":"info","ts":...,"msg":"example log message","exampleKey":"exampleValue"}
}
