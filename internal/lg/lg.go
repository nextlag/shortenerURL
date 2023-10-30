package lg

import (
	"go.uber.org/zap"
)

func New() *zap.Logger {
	// Настраиваем конфигурацию логгера
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel) // Уровень логирования

	// Создаем логгер
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // Отложенное закрытие логгера
	return logger
}
