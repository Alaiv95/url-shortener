package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"urlShortener/internal/config"
	"urlShortener/internal/lib/base62"
	"urlShortener/internal/lib/numGen"
	"urlShortener/internal/storage"
)

type API struct {
	Router *mux.Router
	cfg    *config.HttpServer
	log    *slog.Logger
	db     storage.Storage
}

type SaveUrlReq struct {
	Url string `json:"url" validate:"required,url"`
}

type UrlResp struct {
	Response
	Url string `json:"url"`
}

func New(db storage.Storage, cfg *config.HttpServer, log *slog.Logger) *API {
	api := &API{
		Router: mux.NewRouter(),
		db:     db,
		cfg:    cfg,
		log:    log,
	}

	api.Endpoints()

	return api
}

func (a *API) Endpoints() {
	a.Router.HandleFunc("/api/v1/url", a.Save).Methods("POST")
	a.Router.HandleFunc("/api/v1/url/original", a.Original).Methods("GET")
}

func (a *API) Original(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	shortUrl := r.URL.Query().Get("url")

	origUrl, err := a.db.OrigUrl(shortUrl)

	if err != nil {
		a.log.Error("Error getting url", "err", err.Error())
		WriteResp(w, Err(err.Error()), http.StatusBadRequest)
		return
	}

	responseOk(w, origUrl, http.StatusOK)
}

func (a *API) Save(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req SaveUrlReq

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		a.log.Error("Error decoding json")
		WriteResp(w, Err(err.Error()), http.StatusBadRequest)
		return
	}

	slug := base62.ConvertNum(numGen.Generate())
	shortUrl := fmt.Sprintf("%s/%s", a.cfg.Address, slug)

	resp, err := a.db.SaveUrl(req.Url, shortUrl)
	if err != nil {
		a.log.Error("Error saving link", "err", err.Error())
		WriteResp(w, Err(err.Error()), http.StatusBadRequest)
		return
	}

	responseOk(w, resp, http.StatusOK)
}

func responseOk(w http.ResponseWriter, url string, status int) {
	r := UrlResp{
		Response: Ok(),
		Url:      url,
	}

	WriteResp(w, r, status)
}
