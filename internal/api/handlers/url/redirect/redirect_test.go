package redirect_test

import (
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"urlShortener/internal/api/handlers/url/redirect"
	"urlShortener/internal/api/handlers/url/save"
	"urlShortener/internal/config"
	"urlShortener/internal/storage/memdb"
)

var log *slog.Logger
var getter redirect.UrlGetter
var saver save.UrlSaver

func TestMain(m *testing.M) {
	cfg := config.Config{
		Env: "test",
		Http: config.HttpServer{
			Address:     "test",
			Timeout:     0,
			IdleTimeout: 0,
		},
	}

	log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	db := memdb.New(&cfg)

	getter = db
	saver = db

	os.Exit(m.Run())
}

func TestAPI_Original(t *testing.T) {
	const uri = "http://yandex.ru"
	const shortUrl = "asd"

	_, err := saver.SaveUrl(uri, shortUrl)
	if err != nil {
		t.Errorf("ошибка при подготовке данных")
	}

	h := redirect.New(log, getter)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/url/{slug}", nil)
	req = mux.SetURLVars(req, map[string]string{
		"slug": shortUrl,
	})

	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusFound {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusFound)
	}
}
