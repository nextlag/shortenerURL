package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func New(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			requestLogger := logger.With(
				zap.String("component", "middleware/logger"),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
			)

			ww := NewWrapResponseWriter(w)
			t1 := time.Now()

			defer func() {
				requestLogger.Info("request completed",
					zap.Int("status", ww.Status()),
					zap.Int("bytes", ww.BytesWritten()),
					zap.String("duration", time.Since(t1).String()),
				)
			}()

			next.ServeHTTP(&ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

// NewWrapResponseWriter создает WrapResponseWriter
func NewWrapResponseWriter(w http.ResponseWriter) WrapResponseWriter {
	return WrapResponseWriter{
		ResponseWriter: w,
		status:         http.StatusOK,
		written:        0,
	}
}

// WrapResponseWriter обеспечивает логирование статуса ответа
type WrapResponseWriter struct {
	http.ResponseWriter
	status  int
	written int
}

// WriteHeader сохраняет статус ответа
func (ww *WrapResponseWriter) WriteHeader(code int) {
	ww.status = code
	ww.ResponseWriter.WriteHeader(code)
}

// Write сохраняет количество байтов ответа
func (ww *WrapResponseWriter) Write(b []byte) (int, error) {
	n, err := ww.ResponseWriter.Write(b)
	ww.written += n
	return n, err
}

// Status возвращает статус ответа
func (ww *WrapResponseWriter) Status() int {
	return ww.status
}

// BytesWritten возвращает количество записанных байтов
func (ww *WrapResponseWriter) BytesWritten() int {
	return ww.written
}
