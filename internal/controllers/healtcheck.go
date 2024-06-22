// Package controllers provides HTTP handlers for URL shortening service.
package controllers

import (
	"net/http"

	"go.uber.org/zap"
)

// HealthCheck handles HTTP requests to perform a health check on the service.
func (c *Controller) HealthCheck(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ok, err := c.uc.DoHealthcheck()

		if err != nil {
			c.log.Error("Error performing health check", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if ok {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
