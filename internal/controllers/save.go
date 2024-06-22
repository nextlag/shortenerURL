// Package controllers provides HTTP handlers for URL shortening service.
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
func (c *Controller) Save(w http.ResponseWriter, r *http.Request) {
	// Read the request body (original URL)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request 400", http.StatusBadRequest)
		return
	}

	// Check the authentication cookie
	uuid, err := auth.CheckCookie(w, r, c.log)
	if err != nil {
		c.log.Error("Error getting cookie: ", zap.Error(err))
		return
	}

	// Attempt to save the short URL and original URL in the storage
	alias, err := c.uc.DoPut(r.Context(), string(body), uuid)

	// Handle duplicate URL conflict
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

	if err != nil {
		c.log.Error("failed to add URL", zap.Error(err), zap.String("path to file storage", c.cfg.FileStorage))
		http.Error(w, fmt.Sprintf("failed to add URL: %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	// Send the short URL in the HTTP response body
	_, err = fmt.Fprintf(w, "%s/%s", c.cfg.BaseURL, alias)
	if err != nil {
		c.log.Error("error sending short URL response", zap.Error(err))
		return
	}
	c.log.Info("alias added successfully", zap.String("alias", alias))
}
