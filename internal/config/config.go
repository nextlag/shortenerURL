package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

// ConfigHTTP структура для хранения конфигурации HTTP-сервера.
type ConfigHTTP struct {
	Address     string `json:"address" env:"SERVER_ADDRESS" envDefault:":8080"`
	URLShort    string `json:"url_short" env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStorage string `json:"file_storage,omitempty" env:"FILE_STORAGE_PATH" envDefault:"/tmp/data.json"`
	DSN         string `json:"dsn,omitempty" env:"DATABASE_DSN" envDefault:""`
}

// Config - переменная с конфигурацией
var Config ConfigHTTP

// MakeConfig инициализирует конфигурацию, считывая флаги командной строки и переменные окружения.
func MakeConfig() error {
	// Определение флагов командной строки для настройки конфигурации.
	flag.StringVar(&Config.Address, "a", Config.Address, "Address HTTP-server")
	flag.StringVar(&Config.URLShort, "b", Config.URLShort, "Base URL")
	flag.StringVar(&Config.FileStorage, "f", Config.FileStorage, "Storage in data.json")
	flag.StringVar(&Config.DSN, "d", Config.DSN, "Connect to database")
	return env.Parse(&Config)
}
