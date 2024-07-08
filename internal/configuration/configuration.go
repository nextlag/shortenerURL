package configuration

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/caarlos0/env/v6"
)

var (
	once sync.Once
	cfg  Config
)

// Config structure for configuration.
type Config struct {
	ServerHTTP
	ConfigPath string `json:"config_path" env:"CONFIG_PATH" envDefault:"configuration.json"`
}

// ServerHTTP - structure for storing HTTP server configuration.
type ServerHTTP struct {
	Host          string `json:"host" env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL       string `json:"base_url" env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStorage   string `json:"file_storage,omitempty" env:"FILE_STORAGE_PATH" envDefault:""`
	DSN           string `json:"dsn,omitempty" env:"DATABASE_DSN" envDefault:""`
	EnableHTTPS   bool   `json:"enable_https" env:"ENABLE_HTTPS" envDefault:"false"`
	Cert          string `json:"cert" env:"CERT" envDefault:"cert.pem"`
	Key           string `json:"key" env:"KEY" envDefault:"key.pem"`
	TrustedSubnet string `json:"trusted_subnet" envDefault:""`
}

// Load initializes the configuration by reading command line flags and environment variables.
func Load() (*Config, error) {
	var err error
	once.Do(func() {
		// Регистрируем флаги командной строки
		flag.StringVar(&cfg.Host, "a", cfg.Host, "Host HTTP-server")
		flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "Base URL")
		flag.StringVar(&cfg.FileStorage, "f", cfg.FileStorage, "Storage in data.json")
		flag.StringVar(&cfg.DSN, "d", cfg.DSN, "Connect to database")
		flag.StringVar(&cfg.ConfigPath, "c", cfg.ConfigPath, "Config name file")
		flag.BoolVar(&cfg.EnableHTTPS, "s", cfg.EnableHTTPS, "enabling HTTPS connection")
		flag.StringVar(&cfg.TrustedSubnet, "t", cfg.TrustedSubnet, "trusted subnet address")

		// Получаем путь к конфигурационному файлу из переменных окружения, если указан
		if configPath := os.Getenv("CONFIG_PATH"); configPath != "" {
			cfg.ConfigPath = configPath
			err = loadConfigFromJSON()
			if err != nil {
				return
			}
		}

		if certPath := os.Getenv("CERT"); certPath != "" {
			cfg.Cert = certPath
		}
		if keyPath := os.Getenv("KEY"); keyPath != "" {
			cfg.Key = keyPath
		}

		// override with environment variables
		err = env.Parse(&cfg)
		if err != nil {
			return
		}

		// override with command-line flags
		flag.Parse()
	})
	return &cfg, err
}

func loadConfigFromJSON() error {
	if cfg.ConfigPath == "" {
		log.Println("the path to the configuration file is empty")
		return nil
	}

	data, err := os.ReadFile(cfg.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if err = json.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("failed to parse configuration: %w", err)
	}
	return nil
}
