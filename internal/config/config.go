package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "shorten"
	password = "skypass12345"
	dbname   = "shorten"
)

type PSQLConfig struct {
	Host string `json:"host"`
	Port int    `json:"port,omitempty"`
	User string `json:"user"`
	Pass string `json:"pass"`
	DB   string `json:"db_name"`
}

func NewConnectPSQL() *PSQLConfig {
	return &PSQLConfig{Host: host, Port: port, User: user, Pass: password, DB: dbname}
}

var r = NewConnectPSQL()

func (p *PSQLConfig) ConnectDB() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		r.Host, r.Port, r.User, r.Pass, r.DB)
}

// ArgsHTTP структура для хранения конфигурации HTTP-сервера.
type ArgsHTTP struct {
	Address     string `env:"SERVER_ADDRESS" envDefault:":8080"`
	URLShort    string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStorage string `env:"FILE_STORAGE_PATH" envDefault:"/tmp/data.json"`
	Psql        string `env:"DATABASE_DSN"`
}

// Args - переменная с конфигурацией
var Args ArgsHTTP

// InitializeArgs инициализирует конфигурацию, считывая флаги командной строки и переменные окружения.
func InitializeArgs() error {
	// Определение флагов командной строки для настройки конфигурации.
	flag.StringVar(&Args.Address, "a", Args.Address, "Address HTTP-server")
	flag.StringVar(&Args.URLShort, "b", Args.URLShort, "Base URL")
	flag.StringVar(&Args.FileStorage, "f", Args.FileStorage, "Storage in data.json")
	flag.StringVar(&Args.Psql, "d", Args.Psql, "connect to database")

	// Считывание значений флагов командной строки и переменных окружения в структуру Args.
	if err := env.Parse(&Args); err != nil {
		return err
	}
	Args.Psql = r.ConnectDB()
	return nil
}
