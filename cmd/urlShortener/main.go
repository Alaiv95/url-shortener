package main

import (
	"log/slog"
	"os"
	"urlShortener/internal/config"
	"urlShortener/internal/server"
	"urlShortener/internal/storage/memdb"
)

const (
	local = "local"
	prod  = "prod"
)

func main() {
	cfg := config.MustLoad()
	logger := setupLogger(cfg.Env)

	logger.Debug("Config loaded and Logger enabled")

	storage := memdb.New(cfg)

	logger.Debug("Storage initialized")

	serv := server.New(storage, &cfg.Http, logger)
	serv.Start()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case local:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case prod:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	return log
}
