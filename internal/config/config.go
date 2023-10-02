package config

import (
	"flag"
)

// Cfg - глобальная переменная, хранящая загруженную конфигурацию.
var Cfg = Load()

// Config Структура, представляющая конфигурацию приложения.
type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	HTTPServer `yaml:"http_server"`
}

// HTTPServer Структура, представляющая конфигурацию HTTP-сервера.
type HTTPServer struct {
	Address  string `yaml:"addr" json:"addr,omitempty"`
	URLShort string `yaml:"url_short" json:"url_short,omitempty"`
}

// ParseFlags Парсинг флагов командной строки.
func ParseFlags() (*string, *string) {
	configStartAddr := flag.String("a", ":8080", "Адрес HTTP-сервера")
	configShortURL := flag.String("b", "localhost:8080", "URL короткой ссылки")
	flag.Parse()
	return configStartAddr, configShortURL
}

func Load() *Config {
	configStartAddr, configShortURL := ParseFlags()

	cfg := &Config{
		HTTPServer: HTTPServer{
			Address:  *configStartAddr,
			URLShort: *configShortURL,
		},
	}
	return cfg
}
