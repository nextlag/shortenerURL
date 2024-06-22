// Package logger provides middleware for logging HTTP requests and responses.
package logger

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
)

// RequestFields contains fields for logging HTTP requests.
type RequestFields struct {
	Method         string `json:"method"`                          // HTTP method (GET, POST, etc.)
	Path           string `json:"path"`                            // URL path
	RemoteAddr     string `json:"remote_addr"`                     // Remote address of the client
	UserAgent      string `json:"user_agent"`                      // User agent string
	RequestID      string `json:"request_id"`                      // Request ID for tracing
	DataStorageLoc string `json:"data_storage_location,omitempty"` // Location of data storage
	ContentType    string `json:"content_type,omitempty"`          // Content type of the request
	Status         int    `json:"status"`                          // HTTP status code of the response
	Bytes          int    `json:"bytes,omitempty"`                 // Number of bytes in the response
	Duration       string `json:"duration"`                        // Duration of the request
}

// New creates and returns a new middleware for logging HTTP requests.
func New(log *zap.Logger, cfg config.HTTPServer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			requestFields := RequestFields{
				Method:         r.Method,
				Path:           r.URL.Path,
				RemoteAddr:     r.RemoteAddr,
				UserAgent:      r.UserAgent(),
				RequestID:      middleware.GetReqID(r.Context()),
				ContentType:    r.Header.Get("Content-Type"),
				DataStorageLoc: cfg.FileStorage,
			}

			// Create a WrapResponseWriter to capture the status and byte count of the response.
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now() // Start the timer to measure request duration.

			defer func() {
				// After the request is handled, log the request information including status, byte count, and duration.
				requestFields.Status = ww.Status()
				requestFields.Bytes = ww.BytesWritten()
				requestFields.Duration = time.Since(t1).String()

				// Log only if the request resulted in an error status.
				if requestFields.Status >= http.StatusInternalServerError {
					log.Error("request error: ", zap.Any("error", requestFields))
				} else {
					log.Info("mw", zap.Any("request", requestFields))
				}
			}()

			// Pass the request to the next handler.
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

// SetupLogger initializes and returns a new zap.Logger instance.
func SetupLogger() *zap.Logger {
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel) // Set the logging level

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // Defer the logger sync
	return logger
}
