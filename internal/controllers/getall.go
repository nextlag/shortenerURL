// Package controllers provides HTTP handlers for URL shortening service.
package controllers

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/usecase/auth"
)

// GetAll handles HTTP requests to retrieve all URLs associated with a user.
func (c *Controller) GetAll(w http.ResponseWriter, r *http.Request) {
	// Check the authentication cookie.
	userID, err := auth.CheckCookie(w, r, c.log)
	if err != nil {
		c.log.Error("Unauthorized access: ", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

	// Retrieve all URLs associated with the user.
	userURLs, err := c.uc.DoGetAll(r.Context(), userID, c.cfg.BaseURL)
	if err != nil {
		c.log.Error("Error getting URLs by ID", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	// Set the response headers and write the response.
	w.Header().Set("Content-Type", "application/json")
	if string(userURLs) == "null" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("No content"))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(userURLs)
	}
}
