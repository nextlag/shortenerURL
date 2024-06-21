package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nextlag/shortenerURL/internal/middleware/logger"
)

func TestCheckCookie(t *testing.T) {
	tests := []struct {
		name         string
		cookieValue  string
		expectedCode int
	}{
		{
			name:         "Success",
			cookieValue:  "",
			expectedCode: http.StatusOK,
		},
		{
			name:         "CookiePresent",
			cookieValue:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo5MjR9.3-u9JqljBCsqfT_yrPfIpwwvKtoryZYuLB_UouDB5TM",
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := logger.SetupLogger()
			req := httptest.NewRequest("GET", "/example", nil)
			rr := httptest.NewRecorder()

			if tt.cookieValue != "" {
				cookie := &http.Cookie{Name: "UserID", Value: tt.cookieValue}
				req.AddCookie(cookie)
			}

			id, err := CheckCookie(rr, req, log)

			assert.NoError(t, err, "Expected no error")
			assert.Equal(t, tt.expectedCode, rr.Code, "Expected status code %d, but got %d", tt.expectedCode, rr.Code)

			if tt.expectedCode == http.StatusOK {
				assert.Greater(t, id, 0, "Expected positive user ID")
			}
		})
	}
}
