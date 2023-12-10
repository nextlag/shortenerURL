package entity

import "time"

type Request struct {
	UserID     int       `json:"user_id"`
	UUID       string    `json:"uuid"`
	Alias      string    `json:"short_url,omitempty"`
	URL       string    `json:"original_url" validate:"required,url"`
	CreatedAt time.Time `json:"created_at"`
}
