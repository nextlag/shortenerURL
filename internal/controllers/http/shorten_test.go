package http

import (
	"encoding/json"
	"log"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/configuration"
)

func TestShorten(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name         string
		body         string
		expectedJSON string
	}{
		{
			name:         "ValidRequest",
			body:         `{"url": "http://example.com", "alias": "example"}`,
			expectedJSON: `{"result":"http://localhost:8080/example"}`,
		},
		{
			name:         "Empty Request Body1",
			body:         `{}`,
			expectedJSON: `{"error":"поле URL обязательно для заполнения"}`,
		},
		{
			name:         "Empty Request Body2",
			body:         `{"url": "example.com"}`,
			expectedJSON: `{"error":"поле URL не является допустимым URL"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg, err := configuration.Load()
			if err != nil {
				log.Fatal("Failed to get configuration")
				return
			}

			_, db, _ := Ctrl(t)
			if !strings.Contains(test.name, "ValidRequest") {
				db.EXPECT().DoPut(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			} else {
				db.EXPECT().DoPut(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("example", nil).Times(1)
			}
			log := zap.NewNop()
			reqBody := strings.NewReader(test.body)
			req := httptest.NewRequest("POST", "/api/shorten", reqBody)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			wg := sync.WaitGroup{}
			New(db, &wg, cfg, log).Shorten(w, req)
			resp := w.Result()
			defer resp.Body.Close()

			var responseJSON map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseJSON)
			require.NoError(t, err)

			var expectedJSON map[string]interface{}
			err = json.NewDecoder(strings.NewReader(test.expectedJSON)).Decode(&expectedJSON)
			require.NoError(t, err)

			assert.Equal(t, expectedJSON, responseJSON)
		})
	}
}
