package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/nextlag/shortenerURL/internal/config"
)

// Response представляет JSON-ответ с полем для сообщения об ошибке.
type Response struct {
	Error  string `json:"error,omitempty"`
	Result string `json:"result"` // Результат - сокращенная ссылка.
}

// Error создает JSON-ответ с заданным сообщением об ошибке.
func Error(msg string) Response {
	return Response{
		Error: msg,
	}
}

// ValidationError создает JSON-ответ с сообщениями об ошибках валидации.
func ValidationError(errs validator.ValidationErrors) Response {
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

// ResponseConflict обрабатывает конфликтный запрос
func ResponseConflict(w http.ResponseWriter, alias string) {
	response := Response{
		Result: fmt.Sprintf("%s/%s", config.Cfg.BaseURL, alias),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusConflict)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Обработка ошибки кодирования JSON.
		http.Error(w, "failed to encode JSON response", http.StatusInternalServerError)
	}
}

// ResponseCreated отправляет успешный ответ с сокращенной ссылкой в JSON-формате.
func ResponseCreated(w http.ResponseWriter, alias string) {
	response := Response{
		Result: fmt.Sprintf("%s/%s", config.Cfg.BaseURL, alias),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Обработка ошибки кодирования JSON.
		http.Error(w, "failed to encode JSON response", http.StatusInternalServerError)
	}
}
