package response

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

// Response представляет JSON-ответ с полем для сообщения об ошибке.
type Response struct {
	Error string `json:"error,omitempty"`
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
			errMsgs = append(errMsgs, fmt.Sprintf("поле %s не является допустимым Url", err.Field()))
		}
	}

	return Response{
		Error: strings.Join(errMsgs, ", "),
	}
}
