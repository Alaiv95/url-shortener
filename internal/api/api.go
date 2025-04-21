package api

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"urlShortener/internal/api/mw"
	"urlShortener/internal/config"
	"urlShortener/internal/kafka"
	"urlShortener/internal/redis"
	"urlShortener/internal/service"
	"urlShortener/internal/service/url"
	"urlShortener/internal/storage"
)

type API struct {
	Router  *mux.Router
	cfg     *config.HttpServer
	log     *slog.Logger
	db      storage.UrlRepo
	kf      *kafka.Client
	cache   *redis.Client
	urlSrvc service.UrlService
}

// UrlReq структура запроса для сохранения новой ссылки
type UrlReq struct {
	Url string `json:"url"`
}

// Resp структура ответа для хендлеров url
type Resp struct {
	Response
	Url string `json:"url"`
}

// UrlSaver интерфейс для сохранения ссылок
type UrlSaver interface {
	SaveUrl(origUrl string, shortUrl string) (string, error)
}

// New конструктор для инициализации Api со всеми зависимостями
func New(
	db storage.UrlRepo,
	cfg *config.HttpServer,
	log *slog.Logger,
	kf *kafka.Client,
	cache *redis.Client) *API {
	a := &API{
		Router:  mux.NewRouter(),
		db:      db,
		cfg:     cfg,
		log:     log,
		kf:      kf,
		cache:   cache,
		urlSrvc: url.New(db, log, kf, cache, cache),
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
	a.Router.HandleFunc("/api/v1/url", a.saveUrl).Methods("POST")
	a.Router.HandleFunc("/api/v1/url/{slug}", a.redirect).Methods("GET")
}

func (a *API) saveUrl(w http.ResponseWriter, r *http.Request) {
	var req UrlReq

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		a.log.Error("Error decoding json")
		WriteRespJson(w, Err(err.Error()), http.StatusBadRequest)
		return
	}

	res, err := a.urlSrvc.SaveUrl(req.Url, a.cfg.Address, context.Background())
	if err != nil {
		a.log.Error("Error saving new url", "error", err)
		WriteRespJson(w, Err(err.Error()), http.StatusBadRequest)
		return
	}

	ResponseOk(w, res, http.StatusOK)

	a.log.Info("Saved new short url: ", "original", req.Url, "short", res)
}

func (a *API) redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortUrl := vars["slug"]

	origUrl, err := a.urlSrvc.Url(shortUrl, context.Background())
	if err != nil {
		a.log.Error("Error getting url", "error", err)
		WriteRespJson(w, Err(err.Error()), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, origUrl, http.StatusFound)

	a.log.Info("Redirect user: ", "short", shortUrl, "original", origUrl)
}
