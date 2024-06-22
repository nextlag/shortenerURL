// Package entity provides data structures for the URL shortening service.
package entity

import "time"

// DBStorage represents the structure for storing URL data in the database.
type DBStorage struct {
	UUID        int       `json:"user_id,omitempty"`      // User ID associated with the URL.
	URL         string    `json:"original_url,omitempty"` // Original URL.
	Alias       string    `json:"short_url,omitempty"`    // Shortened URL alias.
	DeletedFlag bool      `json:"deleted_flag"`           // Flag indicating if the URL is deleted.
	CreatedAt   time.Time `json:"created_at,omitempty"`   // Timestamp when the URL was created.
}
