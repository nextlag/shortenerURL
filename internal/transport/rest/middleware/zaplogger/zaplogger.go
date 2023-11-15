package zaplogger

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
)

// RequestFields содержит поля запроса для логгера.
type RequestFields struct {
	Method         string `json:"method"`
	Path           string `json:"path"`
	RemoteAddr     string `json:"remote_addr"`
	UserAgent      string `json:"user_agent"`
	RequestID      string `json:"request_id"`
	DataStorageLoc string `json:"data_storage_location,omitempty"`
	ContentType    string `json:"content_type,omitempty"`
	Status         int    `json:"status"`
	Bytes          int    `json:"bytes,omitempty"`
	Duration       string `json:"duration"`
}

// New создает и возвращает новый middleware для логирования HTTP запросов.
func New(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// Создаем логгер запроса
			requestFields := RequestFields{
				Method:         r.Method,
				Path:           r.URL.Path,
				RemoteAddr:     r.RemoteAddr,
				UserAgent:      r.UserAgent(),
				RequestID:      middleware.GetReqID(r.Context()),
				ContentType:    r.Header.Get("Content-Type"),
				DataStorageLoc: config.Config.FileStorage,
			}

			// Создаем WrapResponseWriter для перехвата статуса ответа и количества байтов.
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now() // Запускаем таймер для измерения продолжительности запроса.

			defer func() {
				// После выполнения запроса, логируем информацию о запросе, включая статус ответа, количество байтов и продолжительность.
				requestFields.Status = ww.Status()
				requestFields.Bytes = ww.BytesWritten()
				requestFields.Duration = time.Since(t1).String()

				// Добавляем логирование, только если статус запроса - ошибка
				if requestFields.Status >= http.StatusInternalServerError {
					logger.Error("request completed with error", zap.Any("request", requestFields))
				} else {
					logger.Info("", zap.Any("request", requestFields))
				}
			}()

			// Передаем запрос следующему обработчику.
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

// // New создает и возвращает новый middleware для логирования HTTP запросов.
// func New(logger *zap.Logger) func(next http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		fn := func(w http.ResponseWriter, r *http.Request) {
// 			// Создаем логгер запроса
// 			requestLogger := logger.With(
// 				zap.String("method", r.Method),
// 				zap.String("path", r.URL.Path),
// 				zap.String("remote_addr", r.RemoteAddr),
// 				zap.String("user_agent", r.UserAgent()),
// 				zap.String("request_id", middleware.GetReqID(r.Context())),
// 				zap.String("data_storage_location", config.Config.FileStorage),
// 			)
//
// 			contentType := r.Header.Get("Content-Type")
//
// 			if contentType != "" {
// 				requestLogger = requestLogger.With(zap.String("content_type", contentType))
// 			}
//
// 			// Создаем WrapResponseWriter для перехвата статуса ответа и количества байтов.
// 			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
// 			t1 := time.Now() // Запускаем таймер для измерения продолжительности запроса.
//
// 			defer func() {
// 				// После выполнения запроса, логируем информацию о запросе, включая статус ответа, количество байтов и продолжительность.
// 				requestLogger.Info("request completed",
// 					zap.Int("status", ww.Status()),
// 					zap.Int("bytes", ww.BytesWritten()),
// 					zap.String("duration", time.Since(t1).String()),
// 				)
// 			}()
//
// 			// Передаем запрос следующему обработчику.
// 			next.ServeHTTP(ww, r)
// 		}
// 		return http.HandlerFunc(fn)
// 	}
// }
