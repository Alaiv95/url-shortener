package main

import (
	"log/slog"
	"os"
	"urlShortener/internal/config"
	"urlShortener/internal/kafka"
	"urlShortener/internal/redis"
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

	kf, err := kafka.New(cfg.Kafka.Address, cfg.Kafka.UrlTopic, logger)
	if err != nil {
		logger.Error("error initializing Kafka client")
		os.Exit(1)
	}

	logger.Debug("Kafka initialized")

	cache := redis.New(&cfg.Redis)

	logger.Debug("Redis initialized")

	serv := server.New(storage, &cfg.Http, logger, kf, cache)
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
