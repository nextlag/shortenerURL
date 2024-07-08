// Package controllers provides the handlers for managing URL shortening operations.
package controllers

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/usecase"
	"github.com/nextlag/shortenerURL/internal/usecase/auth"
)

// Shorten handles HTTP requests for shortening URLs.
// It decodes the JSON request body, validates the input, checks the user's authentication,
// attempts to save the URL in the storage, and returns the shortened URL or appropriate error messages.
func (c *Controller) Shorten(w http.ResponseWriter, r *http.Request) {
	var req usecase.FileStorage
	// Decode the JSON request body into the Data structure
	err := render.DecodeJSON(r.Body, &req)
	c.log.Info("body", zap.String("URL", req.URL))

	// Handle empty request body
	if errors.Is(err, io.EOF) {
		c.log.Error("request body is empty", zap.Error(err))
		render.JSON(w, r, Error("empty request"))
		return
	}

	// Handle request body decoding errors
	if err != nil {
		c.log.Error("failed to decode request body", zap.Error(err))
		render.JSON(w, r, Error("failed to decode request"))
		return
	}

	// Validate input data using the validator library
	if err = validator.New().Struct(req); err != nil {
		var validateErr validator.ValidationErrors
		errors.As(err, &validateErr)
		c.log.Error("invalid request", zap.Error(err))
		render.JSON(w, r, validationError(validateErr))
		return
	}

	uuid, err := auth.CheckCookie(w, r, c.log)
	if err != nil {
		c.log.Error("Error getting cookie: ", zap.Error(err))
		return
	}

	// Add URL to storage and get the identifier (alias)
	alias, err := c.uc.DoPut(r.Context(), req.URL, req.Alias, uuid)
	if errors.Is(err, usecase.ErrConflict) {
		// Handle conflict error for duplicate URLs
		c.log.Error("trying to add a duplicate URL", zap.Error(err))
		responseConflict(w, alias, c.cfg)
		return
	}

	// Handle errors when adding URL to storage
	if err != nil {
		er := fmt.Sprintf("failed to add URL: %s", err)
		render.JSON(w, r, Error(er))
		return
	}
	// Send the shortened URL to the client
	responseCreated(w, alias, c.cfg)
}
