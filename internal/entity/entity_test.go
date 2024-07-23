package entity

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDBStorage(t *testing.T) {
	original := DBStorage{
		UUID:        1,
		URL:         "https://example.com",
		Alias:       "example",
		DeletedFlag: true,
		CreatedAt:   time.Now(),
	}

	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Errorf("Error marshalling JSON: %v", err)
	}

	var unmarshalled DBStorage
	err = json.Unmarshal(jsonData, &unmarshalled)
	if err != nil {
		t.Errorf("Error unmarshalling JSON: %v", err)
	}
}
