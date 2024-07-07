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
	once       sync.Once
	configPath string
	cfg        Config
)

func init() {
	// Регистрируем флаги командной строки
	flag.StringVar(&cfg.Host, "a", cfg.Host, "Host HTTP-server")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "Base URL")
	flag.StringVar(&cfg.FileStorage, "f", cfg.FileStorage, "Storage in data.json")
	flag.StringVar(&cfg.DSN, "d", cfg.DSN, "Connect to database")
	flag.StringVar(&cfg.ConfigPath, "c", cfg.ConfigPath, "Config name file")
	flag.BoolVar(&cfg.EnableHTTPS, "s", cfg.EnableHTTPS, "enabling HTTPS connection")

	configPath = os.Getenv("CONFIG_PATH")
}

// Config structure for configuration.
type Config struct {
	ServerHTTP
	ConfigPath string `json:"config_path" env:"CONFIG_PATH" envDefault:"configuration.json"`
}

// ServerHTTP - structure for storing HTTP server configuration.
type ServerHTTP struct {
	Host        string `json:"host" env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL     string `json:"base_url" env:"BASE_URL" envDefault:""`
	FileStorage string `json:"file_storage,omitempty" env:"FILE_STORAGE_PATH" envDefault:""`
	DSN         string `json:"dsn,omitempty" env:"DATABASE_DSN" envDefault:""`
	EnableHTTPS bool   `json:"enable_https" env:"ENABLE_HTTPS" envDefault:""`
}

// Load initializes the configuration by reading command line flags and environment variables.
func Load() (*Config, error) {
	var err error
	once.Do(func() {
		// Load configuration from JSON file if specified
		if configPath != "" {
			err = loadConfigFromJSON()
			if err != nil {
				return
			}
		}

		// Override with environment variables
		err = env.Parse(&cfg)
		if err != nil {
			return
		}

		// Override with command-line flags
		flag.Parse()
	})
	return &cfg, err
}

func loadConfigFromJSON() error {
	if configPath == "" {
		log.Println("the path to the configuration file is empty")
		return nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if err = json.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("failed to parse configuration: %w", err)
	}
	return nil
}
