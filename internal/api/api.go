package api

import (
	"github.com/gorilla/mux"
	"log/slog"
	"urlShortener/internal/api/handlers/url/redirect"
	"urlShortener/internal/api/handlers/url/save"
	"urlShortener/internal/api/mw"
	"urlShortener/internal/config"
	"urlShortener/internal/storage/memdb"
)

type API struct {
	Router *mux.Router
	cfg    *config.HttpServer
	log    *slog.Logger
	db     *memdb.Storage
}

// New конструктор для инициализации Api со всеми зависимостями
func New(db *memdb.Storage, cfg *config.HttpServer, log *slog.Logger) *API {
	a := &API{
		Router: mux.NewRouter(),
		db:     db,
		cfg:    cfg,
		log:    log,
	}

	a.Middlewares()
	a.Endpoints()

	return a
}

// Middlewares подключение всех middleware
func (a *API) Middlewares() {
	a.Router.Use(mw.HeadersMiddleware)
}

// Endpoints подключение всех хендлеров
func (a *API) Endpoints() {
	a.Router.Handle("/api/v1/url", save.New(a.log, a.db)).Methods("POST")
	a.Router.HandleFunc("/api/v1/url/{slug}", redirect.New(a.log, a.db)).Methods("GET")
}
