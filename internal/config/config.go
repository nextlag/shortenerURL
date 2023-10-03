package config

import (
	"flag"
)

// HTTPServer Структура, представляющая конфигурацию HTTP-сервера.
var (
	Address  *string
	URLShort *string
)

func init() {
	Address = flag.String("a", ":8080", "Address HTTP-server")
	URLShort = flag.String("b", "127.0.0.1:8080", "Address HTTP-server")
}
