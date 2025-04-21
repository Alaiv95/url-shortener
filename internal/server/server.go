package server

import (
	"log/slog"
	"net/http"
	"os"
	"urlShortener/internal/api"
	"urlShortener/internal/config"
	"urlShortener/internal/kafka"
	"urlShortener/internal/redis"
	"urlShortener/internal/storage/pg"
)

const (
	local = "local"
	prod  = "prod"
)

// Server структура сервера
type Server struct {
	api *api.API
	cfg *config.HttpServer
	log *slog.Logger
}

// New конструктор иницализации сервера с его зависимостями
func New() *Server {
	cfg := config.MustLoad()
	logger := setupLogger(cfg.Env)

	logger.Debug("Config loaded and Logger enabled")

	storage, err := pg.New(cfg.StoragePath)
	if err != nil {
		logger.Error("error initializing storage")
		os.Exit(1)
	}

	logger.Debug("Storage initialized")

	kf, err := kafka.New(cfg.Kafka.Address, cfg.Kafka.UrlTopic, logger)
	if err != nil {
		logger.Error("error initializing Kafka client")
		os.Exit(1)
	}

	logger.Debug("Kafka initialized")

	cache := redis.New(&cfg.Redis)

	logger.Debug("Redis initialized")

	return &Server{
		api: api.New(storage, &cfg.Http, logger, kf, cache),
		log: logger,
		cfg: &cfg.Http,
	}
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

// Start запуск сервера с настройками из конфига
func (s *Server) Start() {
	server := &http.Server{
		Addr:         s.cfg.Address,
		ReadTimeout:  s.cfg.Timeout,
		WriteTimeout: s.cfg.Timeout,
		IdleTimeout:  s.cfg.IdleTimeout,
		Handler:      s.api.Router,
	}

	s.log.Info("Starting server. Will listen on " + s.cfg.Address + "...")

	err := server.ListenAndServe()
	if err != nil {
		s.log.Error("Error serving server", "error", err.Error())
		os.Exit(1)
	}
}
