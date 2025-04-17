package save

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"urlShortener/internal/api/handlers"
	"urlShortener/internal/api/handlers/url"
	"urlShortener/internal/config"
	"urlShortener/internal/lib/base62"
	"urlShortener/internal/lib/numGen"
)

// UrlReq структура запроса для сохранения новой ссылки
type UrlReq struct {
	Url string `json:"url"`
}

// UrlSaver интерфейс для сохранения ссылок
type UrlSaver interface {
	SaveUrl(origUrl string, shortUrl string) (string, error)
}

// New конструктор для иницализации хендлера сохранения ссылки
func New(log *slog.Logger, saver UrlSaver, cfg *config.HttpServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UrlReq

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error("Error decoding json")
			handlers.WriteRespJson(w, handlers.Err(err.Error()), http.StatusBadRequest)
			return
		}

		slug := base62.ConvertNum(numGen.Generate())

		resp, err := saver.SaveUrl(req.Url, slug)
		if err != nil {
			log.Error("Error saving url", "err", err.Error())
			handlers.WriteRespJson(w, handlers.Err(err.Error()), http.StatusBadRequest)
			return
		}

		url.ResponseOk(w, resp, http.StatusOK)
	}
}
