package config

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
)

func init() {
	configPath = os.Getenv("CONFIG_PATH")
}

// HTTPServer structure for storing HTTP server configuration.
type HTTPServer struct {
	Host        string `json:"host" env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL     string `json:"base_url" env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStorage string `json:"file_storage,omitempty" env:"FILE_STORAGE_PATH" envDefault:""`
	DSN         string `json:"dsn,omitempty" env:"DATABASE_DSN" envDefault:""`
	EnableHTTPS bool   `json:"enable_https" env:"ENABLE_HTTPS" envDefault:"false"`
	ConfigPath  string `json:"config_path" env:"CONFIG_PATH" envDefault:"config.json"`
}

// Cfg - variable with HTTP server configuration
var Cfg HTTPServer

// Load initializes the configuration by reading command line flags and environment variables.
func Load() error {
	once.Do(func() {
		if configPath != "" {
			err := loadConfigFromJSON()
			if err != nil {
				return
			}
		}

		// Определение флагов командной строки для настройки конфигурации.
		flag.StringVar(&Cfg.Host, "a", Cfg.Host, "Host HTTP-server")
		flag.StringVar(&Cfg.BaseURL, "b", Cfg.BaseURL, "Base URL")
		flag.StringVar(&Cfg.FileStorage, "f", Cfg.FileStorage, "Storage in data.json")
		flag.StringVar(&Cfg.DSN, "d", Cfg.DSN, "Connect to database")
		flag.StringVar(&Cfg.ConfigPath, "c", Cfg.ConfigPath, "Config name file")
		flag.BoolVar(&Cfg.EnableHTTPS, "s", Cfg.EnableHTTPS, "enabling HTTPS connection")
	})
	return env.Parse(&Cfg)
}

func loadConfigFromJSON() error {
	if configPath == "" {
		log.Println("the path to the config file is empty")
		return nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if err = json.Unmarshal(data, &Cfg); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}
	return nil
}
