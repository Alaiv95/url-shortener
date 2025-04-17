package original

import (
	"log/slog"
	"net/http"
	"urlShortener/internal/api/handlers"
	"urlShortener/internal/api/handlers/url"
)

// UrlGetter интерфейс для получения ссылки
type UrlGetter interface {
	GetUrl(shortUrl string) (string, error)
}

// New конструктор для создания хендлера на получение оригинального url по короткому
func New(log *slog.Logger, db UrlGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortUrl := r.URL.Query().Get("url")

		origUrl, err := db.GetUrl(shortUrl)

		if err != nil {
			log.Error("Error getting url", "err", err.Error())
			handlers.WriteRespJson(w, handlers.Err(err.Error()), http.StatusBadRequest)
			return
		}

		url.ResponseOk(w, origUrl, http.StatusOK)
	}
}
