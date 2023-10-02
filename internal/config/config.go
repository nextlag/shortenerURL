package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Cfg - глобальная переменная, хранящая загруженную конфигурацию.
var Cfg = MustLoad()

// Config Структура, представляющая конфигурацию приложения.
type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	HTTPServer `yaml:"http_server"`
}

// HTTPServer Структура, представляющая конфигурацию HTTP-сервера.
type HTTPServer struct {
	Address     string        `yaml:"addr" addr-default:"localhost:8088" json:"address,omitempty"`
	URLShort    string        `yaml:"URLShort" url_short:"localhost:8088" json:"url_short,omitempty"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s" json:"timeout,omitempty"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s" json:"idle_timeout,omitempty"`
}

// Парсинг флагов командной строки.
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
	// Получаем текущий рабочий каталог
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	configPath, configStartAddr, configShortURL := ParseFlags()

	if *configPath == "" {
		*configPath = filepath.Join(wd, "config/local.yaml")
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
