// Package logger provides middleware for logging HTTP requests and responses.
// It uses the zap logger for structured logging.
package logger

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/configuration"
)

// RequestFields contains the fields for the logger to log HTTP request details.
type RequestFields struct {
	Method         string `json:"method"`                          // HTTP method of the request
	Path           string `json:"path"`                            // URL path of the request
	RemoteAddr     string `json:"remote_addr"`                     // Remote address of the client making the request
	UserAgent      string `json:"user_agent"`                      // User agent of the client making the request
	RequestID      string `json:"request_id"`                      // Unique request ID for tracing
	DataStorageLoc string `json:"data_storage_location,omitempty"` // Location of data storage
	ContentType    string `json:"content_type,omitempty"`          // Content type of the request
	Status         int    `json:"status"`                          // HTTP status code of the response
	Bytes          int    `json:"bytes,omitempty"`                 // Number of bytes written in the response
	Duration       string `json:"duration"`                        // Duration of the request handling
}

// New creates and returns a new middleware for logging HTTP requests.
func New(log *zap.Logger, cfg *configuration.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// Create request logger
			requestFields := RequestFields{
				Method:         r.Method,
				Path:           r.URL.Path,
				RemoteAddr:     r.RemoteAddr,
				UserAgent:      r.UserAgent(),
				RequestID:      middleware.GetReqID(r.Context()),
				ContentType:    r.Header.Get("Content-Type"),
				DataStorageLoc: cfg.FileStorage,
			}

			// Create WrapResponseWriter to intercept response status and byte count.
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now() // Start timer to measure request duration.

			defer func() {
				// After request is processed, log request information including response status, byte count, and duration.
				requestFields.Status = ww.Status()
				requestFields.Bytes = ww.BytesWritten()
				requestFields.Duration = time.Since(t1).String()

				// Log request only if status is an error.
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

// SetupLogger sets up and returns a new zap logger with development configuration.
func SetupLogger() *zap.Logger {
	// Configure the logger
	cfgLogger := zap.NewDevelopmentConfig()
	cfgLogger.Level = zap.NewAtomicLevelAt(zap.DebugLevel) // Set log level
	cfgLogger.DisableStacktrace = true

	// Create the logger
	logger, err := cfgLogger.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // Defer logger sync to flush logs before exit
	return logger
}
