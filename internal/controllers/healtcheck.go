package controllers

import (
	"net/http"

	"go.uber.org/zap"
)

func (c *Controller) HealthCheck(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Выполняем healthcheck
		ok, err := c.uc.DoHealthcheck()

		if err != nil {
			c.log.Error("Ошибка при выполнении healthcheck", zap.Error(err))
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
