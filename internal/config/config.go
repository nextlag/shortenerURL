package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

// Server структура для хранения конфигурации HTTP-сервера.
type Server struct {
	Host    string
	BaseURL string
}

type Option func(a *Server)

func New(option ...Option) Server {
	s := Server{
		Host:    ":8080",
		BaseURL: "http://localhost:8080",
	}
	for _, fn := range option {
		fn(&s)
	}
	return s
}

func WithHost(host, base string) Option {
	return func(s *Server) {
		s.Host = host
		s.BaseURL = base
	}
}

var Args = New()

// Args2 Можно создать новый сервер с дефолтными параметрами
//var Args2 = New(
//	WithHost(":9090", "https://127.0.0.1:9090"))

// InitializeArgs инициализирует конфигурацию, считывая флаги командной строки и переменные окружения.
func InitializeArgs() error {
	// Определение флагов командной строки для настройки конфигурации.
	flag.StringVar(&Args.Host, "a", Args.Host, "Address HTTP-server")
	flag.StringVar(&Args.BaseURL, "b", Args.BaseURL, "Base URL")

	// Считывание значений флагов командной строки и переменных окружения в структуру Args.
	if err := env.Parse(&Args); err != nil {
		return err
	}
	return nil
}
