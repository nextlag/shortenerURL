// Package controllers provides the handlers for managing URL shortening operations.
package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/nextlag/shortenerURL/internal/configuration"
)

// Response represents a JSON response with a field for error messages and the result as a shortened URL.
type Response struct {
	Error  string `json:"error,omitempty"`
	Result string `json:"result,omitempty"` // Result - shortened URL.
}

// Error creates a JSON response with the given error message.
func Error(msg string) Response {
	return Response{
		Error: msg,
	}
}

// validationError creates a JSON response with validation error messages.
func validationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("поле %s обязательно для заполнения", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("поле %s не является допустимым URL", err.Field()))
		}
	}

	return Response{
		Error: strings.Join(errMsgs, ", "),
	}
}

// responseConflict handles a conflict request and returns the existing shortened URL.
func responseConflict(w http.ResponseWriter, alias string, cfg *configuration.Config) {
	response := Response{
		Result: fmt.Sprintf("%s/%s", cfg.BaseURL, alias),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusConflict)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Handle JSON encoding error.
		http.Error(w, "failed to encode JSON response", http.StatusInternalServerError)
	}
}

// responseCreated sends a successful response with the shortened URL in JSON format.
func responseCreated(w http.ResponseWriter, alias string, cfg *configuration.Config) {
	response := Response{
		Result: fmt.Sprintf("%s/%s", cfg.BaseURL, alias),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Handle JSON encoding error.
		http.Error(w, "failed to encode JSON response", http.StatusInternalServerError)
	}
}
