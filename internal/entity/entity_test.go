package entity

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDBStorage(t *testing.T) {
	original := URL{
		UUID:      1,
		URL:       "https://example.com",
		Alias:     "example",
		IsDeleted: true,
		CreatedAt: time.Now(),
	}

	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Errorf("Error marshalling JSON: %v", err)
	}

	var unmarshalled URL
	err = json.Unmarshal(jsonData, &unmarshalled)
	if err != nil {
		t.Errorf("Error unmarshalling JSON: %v", err)
	}
}
