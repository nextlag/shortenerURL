// Package entity provides definitions for data storage structures
// used within the application. These structures include database
// storage representations and their associated metadata.
package entity

import "time"

// URL represents the storage structure for user data in the database.
// It includes fields for the user's unique identifier (UUID), the original URL,
// the shortened URL alias, a flag indicating if the record is deleted, and
// the creation timestamp.
type URL struct {
	UUID      int       `json:"user_id,omitempty"`      // UUID is the unique identifier for the user
	URL       string    `json:"original_url,omitempty"` // URL is the original URL provided by the user
	Alias     string    `json:"short_url,omitempty"`    // Alias is the shortened URL alias
	IsDeleted bool      `json:"is_deleted,omitempty"`   // IsDeleted indicates if the record is marked as deleted
	CreatedAt time.Time `json:"created_at,omitempty"`   // CreatedAt is the timestamp when the record was created
}
