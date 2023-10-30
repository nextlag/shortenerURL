package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

// ArgsHTTP структура для хранения конфигурации HTTP-сервера.
type ArgsHTTP struct {
	Address     string `env:"SERVER_ADDRESS" envDefault:":8080"`
	URLShort    string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStorage string `env:"FILE_STORAGE_PATH" envDefault:"/tmp/data.json"`
}

// Args - переменная с конфигурацией
var Args ArgsHTTP

// InitializeArgs инициализирует конфигурацию, считывая флаги командной строки и переменные окружения.
func InitializeArgs() error {
	// Определение флагов командной строки для настройки конфигурации.
	flag.StringVar(&Args.Address, "a", Args.Address, "Address HTTP-server")
	flag.StringVar(&Args.URLShort, "b", Args.URLShort, "Base URL")
	flag.StringVar(&Args.FileStorage, "f", Args.FileStorage, "Storage in file.json")

	// Считывание значений флагов командной строки и переменных окружения в структуру Args.
	if err := env.Parse(&Args); err != nil {
		return err
	}
	return nil
}
