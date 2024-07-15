package configuration

import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeConfig(t *testing.T) {

	ast := assert.New(t)

	cfg, err := Load()
	if err != nil {
		log.Fatal("Failed to get configuration")
		return
	}
	ast.NoError(err, "Load should not return an error")

	// Сохраняем оригинальные значения переменных окружения
	originalEnv := os.Environ()
	defer func() {
		// Восстанавливаем оригинальные значения переменных окружения после теста
		os.Clearenv()
		for _, envVar := range originalEnv {
			pair := strings.SplitN(envVar, "=", 2)
			os.Setenv(pair[0], pair[1])
		}
	}()

	// Проверяем, что значения полей структуры Config соответствуют ожидаемым
	expectedHost := ":8080"
	expectedBaseURL := "http://localhost:8080"

	ast.Equal(expectedHost, cfg.Host, "Host should be equal")
	ast.Equal(expectedBaseURL, cfg.BaseURL, "Base URL should be equal")
}
