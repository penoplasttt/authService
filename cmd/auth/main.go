package main

import (
	"log/slog"
	"os"

	"github.com/penoplasttt/authService/internal/config"
)

const (
	envLocal = "local"
	envDev = "dev"
	envProd = "prod"
)

func main() {
	//инициализация конфига
	cfg := config.MustLoad()

	//инициализация логгера
	log := setupLogger(cfg.Env)

	log.Info("starting application")

	//инициализация приложения

	//запуск приложения
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "dev":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	
	return log
}
