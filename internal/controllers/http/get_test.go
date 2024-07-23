package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/nextlag/shortenerURL/internal/entity"
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
			ExpectedStatus:   http.StatusNotFound,
			ExpectedLocation: "",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			ctrl, db, _ := Ctrl(t)

			if test.Name == "Valid ID" {
				db.EXPECT().DoGet(gomock.Any(), gomock.Any()).Return(&entity.URL{URL: "http://example.com", Alias: "example", IsDeleted: false}, nil).Times(1)
			} else {
				db.EXPECT().DoGet(gomock.Any(), gomock.Any()).Return(nil, errors.New("error")).Times(1)
			}

			r := httptest.NewRequest("GET", test.RequestPath, nil)
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(ctrl.Get)
			handler(w, r)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, test.ExpectedStatus, resp.StatusCode)
			location := resp.Header.Get("Location")
			assert.Equal(t, test.ExpectedLocation, location)
		})
	}
}
