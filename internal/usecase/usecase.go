package usecase

import (
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

func NewRequestFile(uuid, alias, url string) RequestEntity {
	return &CustomRequest{
		Request: entity.Request{
			UUID:  uuid,
			Alias: alias,
			URL:   url,
		},
	}
}

func NewRequest(userID int, uuid, alias, url string) RequestEntity {
	return &CustomRequest{
		Request: entity.Request{
			UserID: userID,
			UUID:   uuid,
			Alias:  alias,
			URL:    url,
		},
	}
}
