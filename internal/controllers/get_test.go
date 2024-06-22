package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetHandler(t *testing.T) {
	tests := []struct {
		Name             string
		RequestPath      string
		ExpectedStatus   int
		ExpectedLocation string
	}{
		{
			Name:             "Valid ID",
			RequestPath:      "/example",
			ExpectedStatus:   http.StatusTemporaryRedirect,
			ExpectedLocation: "http://example.com",
		},
		{
			Name:             "Invalid ID",
			RequestPath:      "/nonexistent",
			ExpectedStatus:   http.StatusBadRequest,
			ExpectedLocation: "",
		}}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			ctrl, db, _ := Ctrl(t)

			if test.Name == "Valid ID" {
				db.EXPECT().DoGet(gomock.Any(), gomock.Any()).Return("", false, nil).Times(1)
			} else {
				db.EXPECT().DoGet(gomock.Any(), gomock.Any()).Return("", false, errors.New("error")).Times(1)
			}
			// Создаем фейковый запрос
			r := httptest.NewRequest("GET", test.RequestPath, nil)
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(ctrl.Get)
			handler(w, r)

			// Создаем и вызываем handler для маршрута
			resp := w.Result()

			// Проверяем статус кода
			assert.Equal(t, test.ExpectedStatus, resp.StatusCode)
			// Получаем значение Location
			location := resp.Header.Get("Location")
			assert.Empty(t, location)

			// Закрываем тело HTTP-ответа
			require.NoError(t, resp.Body.Close())
		})
	}
}
