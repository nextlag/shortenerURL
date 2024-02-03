package main

import (
	"errors"
	stdLog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/controllers"
	"github.com/nextlag/shortenerURL/internal/middleware/logger"
	"github.com/nextlag/shortenerURL/internal/usecase"
)

// Время ожидания создания таблицы

func setupServer(router http.Handler) *http.Server {
	// Создание HTTP-сервера с указанным адресом и обработчиком маршрутов
	return &http.Server{
		Addr:    config.Cfg.Host, // Получение адреса из настроек
		Handler: router,
	}
}

func main() {
	if err := config.MakeConfig(); err != nil {
		stdLog.Fatal(err)
	}
	// flag.Parse()
	var (
		log = logger.SetupLogger()
		cfg = config.Cfg
		uc  *usecase.UseCase
	)
	log.Info("initialized flags", zap.String("-a", cfg.Host), zap.String("-b", cfg.BaseURL), zap.String("-f", cfg.FileStorage), zap.String("-d", cfg.DSN))
	// Создание хранилища данных в памяти
	if cfg.FileStorage != "" {
		err := usecase.Load(cfg.FileStorage)
		if err != nil {
			log.Error("failed to load data from file", zap.Error(err))
			os.Exit(1) // Завершаем программу при ошибке загрузки данных
		}
	}

	if cfg.DSN != "" {
		db, err := usecase.NewDB(cfg.DSN, log)
		if err != nil {
			log.Fatal("failed to connect in database", zap.Error(err))
		}
		defer func() {
			if err := db.Stop(); err != nil {
				log.Error("error stopping DB", zap.Error(err))
			}
		}()
		uc = usecase.New(db, log, cfg)
	} else {
		db := usecase.NewData(log, cfg)
		uc = usecase.New(db, log, cfg)
	}

	// Создание и настройка маршрутов и HTTP-сервера
	controller := controllers.New(uc, log, cfg)

	r := chi.NewRouter()
	r.Mount("/", controller.Router(r))

	srv := setupServer(r)

	log.Info("server starting",
		zap.String("address", cfg.Host),
		zap.String("url", cfg.BaseURL))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Запуск HTTP-сервера в горутине
	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			// Если сервер не стартовал, логируем ошибку
			log.Error("failed to start server", zap.String("error", err.Error()))
			done <- os.Interrupt
		}
	}()
	log.Info("server started")

	<-done // Ожидание сигнала завершения
	log.Info("server stopped")
}
