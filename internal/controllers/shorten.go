// Package controllers provides HTTP handlers for URL shortening service.
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
func (c *Controller) Shorten(w http.ResponseWriter, r *http.Request) {
	var req usecase.FileStorage
	err := render.DecodeJSON(r.Body, &req)

	// Handle the case where the request body is empty.
	if errors.Is(err, io.EOF) {
		c.log.Error("request body is empty", zap.Error(err))
		render.JSON(w, r, Error("empty request"))
		return
	}

	// Handle errors decoding the request body.
	if err != nil {
		c.log.Error("failed to decode request body", zap.Error(err))
		render.JSON(w, r, Error("failed to decode request"))
		return
	}

	// Validate the input data using the validator library.
	if err = validator.New().Struct(req); err != nil {
		var validateErr validator.ValidationErrors
		errors.As(err, &validateErr)
		c.log.Error("invalid request", zap.Error(err))
		render.JSON(w, r, ValidationError(validateErr))
		return
	}

	// Check the authentication cookie.
	uuid, err := auth.CheckCookie(w, r, c.log)
	if err != nil {
		c.log.Error("Error getting cookie: ", zap.Error(err))
		return
	}

	// Add the URL to the storage and get the alias.
	alias, err := c.uc.DoPut(r.Context(), req.URL, uuid)
	if errors.Is(err, usecase.ErrConflict) {
		// Handle the case of a duplicate URL conflict.
		c.log.Error("trying to add a duplicate URL", zap.Error(err))
		ResponseConflict(w, alias)
		return
	}

	// Handle errors adding the URL to the storage.
	if err != nil {
		er := fmt.Sprintf("failed to add URL: %s", err)
		render.JSON(w, r, Error(er))
		return
	}

	// Send the response to the client with the shortened URL.
	ResponseCreated(w, alias)
}
