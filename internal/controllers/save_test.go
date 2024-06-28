package controllers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/nextlag/shortenerURL/internal/usecase"
)

func TestSaveHandler(t *testing.T) {
	tests := []struct {
		Name           string
		RequestBody    string
		ExpectedStatus int
	}{
		{
			Name:           "Valid URL",
			RequestBody:    "http://example.com",
			ExpectedStatus: http.StatusCreated,
		},
		{
			Name:           "Duplicate URL",
			RequestBody:    "http://duplicate.com",
			ExpectedStatus: http.StatusConflict,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			ctrl, db, _ := Ctrl(t)

			switch test.Name {
			case "Valid URL":
				db.EXPECT().DoPut(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("newAlias", nil).Times(1)
			case "Duplicate URL":
				db.EXPECT().DoPut(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("duplicateAlias", usecase.ErrConflict).Times(1)
			case "Invalid Request Body":
				// No need to mock db.DoPut for invalid request body case
			}

			// Создаем фейковый запрос
			req := httptest.NewRequest("POST", "/", strings.NewReader(test.RequestBody))
			w := httptest.NewRecorder()

			// Вызываем обработчик
			ctrl.Save(w, req)

			// Проверяем результат
			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, test.ExpectedStatus, resp.StatusCode)
		})
	}
}
