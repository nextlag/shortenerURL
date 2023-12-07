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

func NewRequest(uuid, alias, url string) RequestEntity {
	return &CustomRequest{
		Request: entity.Request{
			UUID:  uuid,
			Alias: alias,
			URL:   url,
		},
	}
}
