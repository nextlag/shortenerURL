package config

import (
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeConfig(t *testing.T) {
	// Используем assert для удобства проверок
	ast := assert.New(t)

	// Вызываем MakeConfig()
	err := MakeConfig()
	ast.NoError(err, "MakeConfig should not return an error")

	// Парсим флаги после вызова MakeConfig()
	flag.Parse()

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

	// Проверяем, что значения полей структуры HTTPServer соответствуют ожидаемым
	expectedHost := ":8080"
	expectedBaseURL := "http://localhost:8080"
	expectedFileStorage := ""
	expectedDSN := ""

	ast.Equal(expectedHost, Cfg.Host, "Host should be equal")
	ast.Equal(expectedBaseURL, Cfg.BaseURL, "Base URL should be equal")
	ast.Equal(expectedFileStorage, Cfg.FileStorage, "File storage should be equal")
	ast.Equal(expectedDSN, Cfg.DSN, "DSN should be equal")
}