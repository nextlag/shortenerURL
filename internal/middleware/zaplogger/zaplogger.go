// Пакет logger предоставляет middleware для логирования HTTP запросов с использованием библиотеки Zap.

package logger

import (
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// New создает и возвращает новый middleware для логирования HTTP запросов.
func New(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// Создаем логгер запроса, добавляя информацию о методе, пути, IP-адресе и User-Agent.
			requestLogger := logger.With(
				zap.String("component", "middleware/zaplogger"),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
			)

			// Создаем WrapResponseWriter для перехвата статуса ответа и количества байтов.
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now() // Запускаем таймер для измерения продолжительности запроса.

			defer func() {
				// После выполнения запроса, логируем информацию о запросе, включая статус ответа, количество байтов и продолжительность.
				requestLogger.Info("request completed",
					zap.Int("status", ww.Status()),
					zap.Int("bytes", ww.BytesWritten()),
					zap.String("duration", time.Since(t1).String()),
				)
			}()

			// Передаем запрос следующему обработчику.
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}