package usecase

import (
	"time"

	"github.com/nextlag/shortenerURL/internal/entity"
)

type RequestEntity interface {
	GetEntityRequest() *entity.Request
}

type CustomRequest struct {
	entity.Request
}

func (r *CustomRequest) GetEntityRequest() *entity.Request {
	return &r.Request
}

func NewRequest(userID int, uuid, alias, url string, created_at time.Time) RequestEntity {
	return &CustomRequest{
		Request: entity.Request{
			UserID:     userID,
			UUID:       uuid,
			Alias:      alias,
			URL:        url,
			Created_at: created_at,
		},
	}
}
