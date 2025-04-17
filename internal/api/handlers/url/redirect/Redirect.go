package redirect

import (
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"urlShortener/internal/api/handlers"
)

// UrlGetter интерфейс для получения ссылки
type UrlGetter interface {
	GetUrl(shortUrl string) (string, error)
}

// New конструктор для создания хендлера для редиректа с короткой ссылки на оригинальную
func New(log *slog.Logger, db UrlGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		shortUrl := vars["slug"]

		origUrl, err := db.GetUrl(shortUrl)

		if err != nil {
			log.Error("Error getting url", "err", err.Error())
			handlers.WriteRespJson(w, handlers.Err(err.Error()), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, origUrl, http.StatusFound)
	}
}
