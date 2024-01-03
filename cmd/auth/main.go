package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/penoplasttt/authService/internal/app"
	"github.com/penoplasttt/authService/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	//инициализация конфига
	cfg := config.MustLoad()

	//инициализация логгера
	log := setupLogger(cfg.Env)

	log.Info("starting application")

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	go application.GRPCSrv.MustRun()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan

	log.Info("stopping application", slog.String("signal",sig.String()))

	application.GRPCSrv.Stop()

	log.Info("application stopped")

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
