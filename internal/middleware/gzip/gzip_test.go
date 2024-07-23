package gzip

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGzipMiddleware(t *testing.T) {
	// Структура с параметрами для тестовых случаев.
	testCases := []struct {
		name              string
		acceptEncoding    string
		expectedBodyEmpty bool
		expectedHeader    string
	}{
		{
			name:              "Accepts_Gzip_Encoding",
			acceptEncoding:    "gzip",
			expectedBodyEmpty: false,
			expectedHeader:    "TestHeader",
		},
		{
			name:              "No_Gzip_Encoding",
			acceptEncoding:    "",
			expectedBodyEmpty: true,
			expectedHeader:    "TestHeader",
		},
	}

	// Создаем тестовый HTTP-обработчик для тестирования middleware.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Custom-Header", "TestHeader")
		w.WriteHeader(http.StatusOK)
	})

	// Создаем экземпляр middleware с тестовым обработчиком.
	middleware := New()

	// Проходим по всем тестовым случаям.
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Создаем тестовый запрос с указанным заголовком Accept-Encoding.
			r := httptest.NewRequest("GET", "http://example.com", nil)
			r.Header.Set("Accept-Encoding", tc.acceptEncoding)

			// Создаем тестовую запись (Recorder) для записи ответа.
			rr := httptest.NewRecorder()

			// Запускаем middleware.
			middleware(handler).ServeHTTP(rr, r)

			// Проверяем статус код ответа.
			assert.Equal(t, http.StatusOK, rr.Code, "Expected status code %d, but got %d", http.StatusOK, rr.Code)

			// Проверяем, что middleware вернул ожидаемое тело ответа.
			assert.Equal(t, tc.expectedBodyEmpty, len(rr.Body.Bytes()) == 0, "Unexpected response body")

			// Проверяем, что middleware добавил заголовок X-Custom-Header.
			assert.Equal(t, tc.expectedHeader, rr.Header().Get("X-Custom-Header"), "Unexpected X-Custom-Header value")
		})
	}
}

func BenchmarkGzipMiddleware(b *testing.B) {
	// Создаем экземпляр middleware.
	middleware := New()

	// Создаем тестовый HTTP-обработчик для бенчмарка.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Запускаем бенчмарк.
	for i := 0; i < b.N; i++ {
		// Создаем тестовый запрос с различными значениями заголовков Accept-Encoding и Content-Encoding.
		r := httptest.NewRequest("GET", "http://example.com", nil)
		r.Header.Set("Accept-Encoding", "gzip")
		r.Header.Set("Content-Encoding", "gzip")

		// Создаем тестовую запись (Recorder) для записи ответа.
		rr := httptest.NewRecorder()

		// Запускаем middleware.
		middleware(handler).ServeHTTP(rr, r)
	}
}
