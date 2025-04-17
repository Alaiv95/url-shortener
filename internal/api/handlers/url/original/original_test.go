package original_test

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"urlShortener/internal/api/handlers/url"
	"urlShortener/internal/api/handlers/url/original"
	"urlShortener/internal/api/handlers/url/save"
	"urlShortener/internal/config"
	"urlShortener/internal/storage/memdb"
)

var log *slog.Logger
var getter original.UrlGetter
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
	const shortUrl = "test/2"

	_, err := saver.SaveUrl(uri, shortUrl)
	if err != nil {
		t.Errorf("ошибка при подготовке данных")
	}

	h := original.New(log, getter)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/url/original", nil)
	q := req.URL.Query()
	q.Add("url", shortUrl)
	req.URL.RawQuery = q.Encode()

	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}

	var resp url.Resp
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Errorf("ошибка при десериализации ответа")
	}

	if uri != resp.Url {
		t.Errorf("url неверен: получили %s, а хотели %s", resp.Url, uri)
	}
}
