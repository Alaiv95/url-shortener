package server

import (
	"net/http"
	"os"
	"urlShortener/internal/api"
)

// Server структура сервера
type Server struct {
	provider *serviceProvider
	api      *api.API
}

// New конструктор иницализации сервера с его зависимостями
func New() *Server {
	server := &Server{
		provider: newServiceProvider(),
	}

	server.InitApi()

	return server
}

// Start запуск сервера с настройками из конфига
func (s *Server) Start() {
	server := &http.Server{
		Addr:         s.provider.Config().Http.Address,
		ReadTimeout:  s.provider.Config().Http.Timeout,
		WriteTimeout: s.provider.Config().Http.Timeout,
		IdleTimeout:  s.provider.Config().Http.IdleTimeout,
		Handler:      s.api.Router,
	}

	s.provider.Logger().Info("Starting server. Will listen on " + s.provider.
		Config().Http.Address + "...")

	err := server.ListenAndServe()
	if err != nil {
		s.provider.Logger().Error("Error serving server", "error", err.Error())
		os.Exit(1)
	}
}

func (s *Server) InitApi() {
	s.api = api.New(
		s.provider.UrlService(),
		&s.provider.Config().Http,
		s.provider.Logger(),
		s.provider.Kafka(),
		s.provider.Redis())
}
