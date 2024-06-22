// Package config provides configuration structures and initialization functions for the HTTP server.
package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

// HTTPServer is a structure for storing HTTP server configuration.
type HTTPServer struct {
	Host        string `json:"host" env:"SERVER_ADDRESS" envDefault:":8080"`                 // Server address.
	BaseURL     string `json:"url_short" env:"BASE_URL" envDefault:"http://localhost:8080"`  // Base URL.
	FileStorage string `json:"file_storage,omitempty" env:"FILE_STORAGE_PATH" envDefault:""` // Path to the data storage file.
	DSN         string `json:"dsn,omitempty" env:"DATABASE_DSN" envDefault:""`               // Database connection string.
}

// Cfg is the variable holding the server configuration.
var Cfg HTTPServer

// MakeConfig initializes the configuration by reading command-line flags and environment variables.
func MakeConfig() error {
	// Define command-line flags for configuration.
	flag.StringVar(&Cfg.Host, "a", Cfg.Host, "Host HTTP-server")                   // Flag for setting the server address.
	flag.StringVar(&Cfg.BaseURL, "b", Cfg.BaseURL, "Base URL")                     // Flag for setting the base URL.
	flag.StringVar(&Cfg.FileStorage, "f", Cfg.FileStorage, "Storage in data.json") // Flag for setting the path to the data storage file.
	flag.StringVar(&Cfg.DSN, "d", Cfg.DSN, "Connect to database")                  // Flag for setting the database connection string.
	return env.Parse(&Cfg)                                                         // Parse environment variables and assign them to the respective fields of the structure.
}
