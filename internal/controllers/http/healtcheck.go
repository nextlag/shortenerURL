// Package controllers provides the handlers for managing URL shortening operations.
package http

import (
	"net/http"

	"go.uber.org/zap"
)

// HealthCheck handles the HTTP request for checking the health of the service.
// It performs a health check by calling the use case layer and responds with
// HTTP 200 OK if the service is healthy, or HTTP 500 Internal Server Error if not.
// If the request method is not GET, it responds with HTTP 405 Method Not Allowed.
func (c *Controller) HealthCheck(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ok, err := c.uc.DoHealthcheck()
		if err != nil {
			c.log.Error("Ошибка при выполнении healthcheck", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if ok {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success!"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
