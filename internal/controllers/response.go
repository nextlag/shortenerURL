// Package controllers provides HTTP handlers for URL shortening service.
package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/nextlag/shortenerURL/internal/config"
)

// Response represents a JSON response with a field for an error message.
type Response struct {
	Error  string `json:"error,omitempty"` // Error message.
	Result string `json:"result"`          // Result message.
}

// Error creates a JSON response with the given error message.
func Error(msg string) Response {
	return Response{
		Error: msg,
	}
}

// ValidationError creates a JSON response with validation error messages.
func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid URL", err.Field()))
		}
	}

	return Response{
		Error: strings.Join(errMsgs, ", "),
	}
}

// ResponseConflict handles a conflicting request by sending a JSON response with the conflicting URL.
func ResponseConflict(w http.ResponseWriter, alias string) {
	response := Response{
		Result: fmt.Sprintf("%s/%s", config.Cfg.BaseURL, alias),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusConflict)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode JSON response", http.StatusInternalServerError)
	}
}

// ResponseCreated sends a successful response with the shortened URL in JSON format.
func ResponseCreated(w http.ResponseWriter, alias string) {
	response := Response{
		Result: fmt.Sprintf("%s/%s", config.Cfg.BaseURL, alias),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode JSON response", http.StatusInternalServerError)
	}
}
