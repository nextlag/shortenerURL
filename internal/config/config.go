package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

// HTTPServer структура для хранения конфигурации HTTP-сервера.
type HTTPServer struct {
	Host        string `json:"host" env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL     string `json:"url_short" env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStorage string `json:"file_storage,omitempty" env:"FILE_STORAGE_PATH" envDefault:""`
	DSN         string `json:"dsn,omitempty" env:"DATABASE_DSN" envDefault:""` // host=localhost port=5432 user=shorten password=skypass12345 database=shorten sslmode=disable
}

// Cfg - переменная с конфигурацией
var Cfg HTTPServer

// MakeConfig инициализирует конфигурацию, считывая флаги командной строки и переменные окружения.
func MakeConfig() error {
	// Определение флагов командной строки для настройки конфигурации.
	flag.StringVar(&Cfg.Host, "a", Cfg.Host, "Host HTTP-server")
	flag.StringVar(&Cfg.BaseURL, "b", Cfg.BaseURL, "Base URL")
	flag.StringVar(&Cfg.FileStorage, "f", Cfg.FileStorage, "Storage in data.json")
	flag.StringVar(&Cfg.DSN, "d", Cfg.DSN, "Connect to database")
	// flag.Parse()
	return env.Parse(&Cfg)
}
