package save

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
	"urlShortener/internal/api/handlers"
	"urlShortener/internal/api/handlers/url"
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

type Producer interface {
	Produce(msgVal []byte, ctx context.Context)
}

type TimedSetter interface {
	Set(key string, v any, time time.Duration) error
}

// New конструктор для иницализации хендлера сохранения ссылки
func New(log *slog.Logger, saver UrlSaver, kf Producer, domain string, cache TimedSetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UrlReq

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error("Error decoding json")
			handlers.WriteRespJson(w, handlers.Err(err.Error()), http.StatusBadRequest)
			return
		}

		log.Info("Start saving short url for:  ", "url", req.Url)

		slug := base62.ConvertNum(numGen.Generate())

		resp, err := saver.SaveUrl(req.Url, slug)
		if err != nil {
			log.Error("Error saving url", "err", err.Error())
			handlers.WriteRespJson(w, handlers.Err(err.Error()), http.StatusBadRequest)
			return
		}

		// form short url with current domain
		shortUrl := fmt.Sprintf("%s/api/v1/url/%s", domain, resp)
		url.ResponseOk(w, shortUrl, http.StatusOK)

		log.Info("Saved new short url: ", "original", req.Url, "short", shortUrl)

		go func() {
			// produce message to kafka that url created
			kf.Produce([]byte(shortUrl), r.Context())

			// set url to cache for 12 hr
			err = cache.Set(slug, req.Url, time.Hour*12)
			if err != nil {
				log.Error("Error saving url to cache", "err", err.Error())
			}
		}()
	}
}
