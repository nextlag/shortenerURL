// Package controllers provides the handlers for managing URL shortening operations.
package http

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/usecase/auth"
	"github.com/nextlag/shortenerURL/internal/usecase/repository/inmemory"
	"github.com/nextlag/shortenerURL/internal/usecase/repository/psql"
)

// Shorten handles HTTP requests for shortening URLs.
// It decodes the JSON request body, validates the input, checks the user's authentication,
// attempts to save the URL in the storage, and returns the shortened URL or appropriate error messages.
func (c *Controller) Shorten(w http.ResponseWriter, r *http.Request) {
	var req inmemory.FileStorage
	err := render.DecodeJSON(r.Body, &req)
	c.log.Info("body", zap.String("URL", req.URL))

	if errors.Is(err, io.EOF) {
		c.log.Error("request body is empty", zap.Error(err))
		render.JSON(w, r, Error("empty request"))
		return
	}

	if err != nil {
		c.log.Error("failed to decode request body", zap.Error(err))
		render.JSON(w, r, Error("failed to decode request"))
		return
	}

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

	alias, err := c.uc.DoPut(r.Context(), req.URL, req.Alias, uuid)
	if errors.Is(err, psql.ErrConflict) {
		c.log.Error("trying to add a duplicate URL", zap.Error(err))
		responseConflict(w, alias, c.cfg)
		return
	}

	if err != nil {
		er := fmt.Sprintf("failed to add URL: %s", err)
		render.JSON(w, r, Error(er))
		return
	}
	responseCreated(w, alias, c.cfg)
}
