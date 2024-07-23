package main

import (
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/nextlag/shortenerURL/internal/configuration"
)

func TestMain(m *testing.M) {
	if _, err := configuration.Load(); err != nil {
		log.Fatal(err)
	}

	exitCode := m.Run()
	// Закрываем все оставшиеся тела ответов
	http.DefaultTransport.(*http.Transport).CloseIdleConnections()
	// Завершаем выполнение программы с кодом завершения.
	os.Exit(exitCode)
}

func TestSetupServer(t *testing.T) {
	// Создаем экземпляр роутера (заменяем реальный роутер на фейковый)
	router := http.NewServeMux()
	cfg, err := configuration.Load()
	if err != nil {
		log.Fatal("Failed to get configuration")
		return
	}

	// Инициализируем сервер
	srv := &http.Server{
		Addr:    cfg.Host,
		Handler: router,
	}

	// Проверяем, что адрес сервера устанавливается правильно
	expectedAddr := cfg.Host
	if srv.Addr != expectedAddr {
		t.Errorf("expected address %s, got %s", expectedAddr, srv.Addr)
	}
	// Проверяем, что обработчик маршрутов устанавливается правильно
	if srv.Handler != router {
		t.Error("expected router handler to be set")
	}
}
