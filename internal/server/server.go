package server

import (
	"log/slog"
	"net/http"
	"os"
	"urlShortener/internal/api"
	"urlShortener/internal/config"
	"urlShortener/internal/storage/memdb"
)

// Server структура сервера
type Server struct {
	api *api.API
	cfg *config.HttpServer
	log *slog.Logger
	db  *memdb.Storage
}

// New конструктор иницализации сервера с его зависимостями
func New(db *memdb.Storage, cfg *config.HttpServer, log *slog.Logger) *Server {
	return &Server{
		db:  db,
		api: api.New(db, cfg, log),
		cfg: cfg,
		log: log,
	}
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
