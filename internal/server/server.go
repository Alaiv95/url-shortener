package server

import (
	"log/slog"
	"net"
	"net/http"
	"os"
	"urlShortener/internal/api"
	"urlShortener/internal/config"
	"urlShortener/internal/storage"
)

type Server struct {
	api *api.API
	cfg *config.HttpServer
	log *slog.Logger
	db  storage.Storage
}

func New(db storage.Storage, cfg *config.HttpServer, log *slog.Logger) *Server {
	return &Server{
		db:  db,
		api: api.New(db, cfg, log),
		cfg: cfg,
		log: log,
	}
}

func (s *Server) Start() {
	server := &http.Server{
		Addr:         s.cfg.Address,
		ReadTimeout:  s.cfg.Timeout,
		WriteTimeout: s.cfg.Timeout,
		IdleTimeout:  s.cfg.IdleTimeout,
		Handler:      s.api.Router,
	}

	lis, err := net.Listen("tcp4", s.cfg.Address)
	if err != nil {
		s.log.Error("Error starting server", "error", err.Error())
		os.Exit(1)
	}

	s.log.Info("Starting server. Will listen on " + s.cfg.Address + "...")

	err = server.Serve(lis)
	if err != nil {
		s.log.Error("Error serving server", "error", err.Error())
		os.Exit(1)
	}
}
