package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

// Cfg - глобальная переменная, хранящая загруженную конфигурацию.
var (
	Cfg        = MustLoad()
	ConfigPath = "../../config/local.yaml"
)

// Config Структура, представляющая конфигурацию приложения.
type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	HTTPServer `yaml:"http_server"`
}

// HTTPServer Структура, представляющая конфигурацию HTTP-сервера.
type HTTPServer struct {
	Address  string `yaml:"addr" addr-default:":8080" json:"addr,omitempty"`
	URLShort string `yaml:"url_short" url_short:"localhost:8080" json:"url_short,omitempty"`
}

// ParseFlags Парсинг флагов командной строки.
func ParseFlags() (*string, *string, *string) {
	configPath := flag.String("config", "", "Путь к конфигурационному файлу")
	configStartAddr := flag.String("a", "", "Адрес HTTP-сервера")
	configShortURL := flag.String("b", "", "URL короткой ссылки")
	flag.Parse()
	return configPath, configStartAddr, configShortURL
}

// Загрузка конфигурационных данных из файла configPath в структуру cfg
func loadConfig(configPath string) (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func MustLoad() *Config {

	configPath, configStartAddr, configShortURL := ParseFlags()

	if *configPath == "" {
		*configPath = ConfigPath
	}
	// Вывод пути к конфигурационному файлу (для проверки)
	fmt.Printf("Путь к конфигурационному файлу: %s\n", *configPath)

	// Проверка существует ли файл по указанному пути configPath
	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", *configPath)
	}

	cfg, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	// задаем значение по умолчанию
	if *configStartAddr != "" {
		cfg.HTTPServer.Address = *configStartAddr
	}
	if *configShortURL != "" {
		cfg.HTTPServer.URLShort = *configShortURL
	}

	return cfg
}
