package entity

import "time"

type DBStorage struct {
	UUID        int       `json:"user_id,omitempty" `
	URL         string    `json:"original_url,omitempty"`
	Alias       string    `json:"short_url,omitempty"`
	DeletedFlag bool      `json:"deleted_flag"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}
