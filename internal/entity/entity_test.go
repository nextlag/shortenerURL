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
		DeletedFlag: false,
		CreatedAt:   time.Now(),
	}

	// Преобразуем структуру в JSON
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Errorf("Error marshalling JSON: %v", err)
	}

	// Демаршализуем JSON обратно в структуру
	var unmarshalled DBStorage
	err = json.Unmarshal(jsonData, &unmarshalled)
	if err != nil {
		t.Errorf("Error unmarshalling JSON: %v", err)
	}
}
