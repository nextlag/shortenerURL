// Package controllers provides the handlers for managing URL shortening operations.
package controllers

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/usecase"
	"github.com/nextlag/shortenerURL/internal/usecase/auth"
)

// Save handles POST requests to create and save a URL in the storage.
// It reads the request body to get the original URL, checks the user's authentication cookie,
// attempts to save the short URL and the original URL in the storage, and handles any conflicts or errors.
func (c *Controller) Save(w http.ResponseWriter, r *http.Request) {

	// Read the request body (original URL)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request 400", http.StatusBadRequest)
		return
	}
	uuid, err := auth.CheckCookie(w, r, c.log)
	if err != nil {
		c.log.Error("Error getting cookie: ", zap.Error(err))
		return
	}

	// Attempt to save the short URL and the original URL in the storage
	alias, err := c.uc.DoPut(r.Context(), string(body), uuid)

	// Handle duplicate conflict
	if errors.Is(err, usecase.ErrConflict) {
		c.log.Error("duplicate url", zap.String("alias", alias), zap.String("url", string(body)))
		w.WriteHeader(http.StatusConflict)
		_, err = fmt.Fprintf(w, "%s/%s", c.cfg.BaseURL, alias)
		if err != nil {
			c.log.Error("error sending short URL response", zap.Error(err))
			return
		}
		return
	}

	// Handle other errors
	if err != nil {
		c.log.Error("failed to add URL", zap.Error(err), zap.String("path to file storage", c.cfg.FileStorage))
		http.Error(w, fmt.Sprintf("failed to add URL: %s", err), http.StatusInternalServerError)
		return
	}
	// Set HTTP status 201 Created
	w.WriteHeader(http.StatusCreated)

	// Send the short URL in the HTTP response body
	_, err = fmt.Fprintf(w, "%s/%s", c.cfg.BaseURL, alias)
	if err != nil {
		c.log.Error("error sending short URL response", zap.Error(err))
		return
	}
	c.log.Info("alias added success", zap.String("alias", alias))
}
