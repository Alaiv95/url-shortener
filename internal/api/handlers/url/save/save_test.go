package save_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"
	"urlShortener/internal/api/handlers/url"
	"urlShortener/internal/api/handlers/url/save"
	"urlShortener/internal/config"
	"urlShortener/internal/storage/memdb"
)

var log *slog.Logger
var saver save.UrlSaver
var cfg config.Config

func TestMain(m *testing.M) {
	cfg = config.Config{
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

	saver = memdb.New(&cfg)

	os.Exit(m.Run())
}

func TestAPI_Save(t *testing.T) {
	data := save.UrlReq{
		Url: "http://google.com",
	}
	payload, _ := json.Marshal(data)
	h := save.New(log, saver, &cfg.Http)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/url", bytes.NewBuffer(payload))

	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}

	var resp url.Resp
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Errorf("ошибка при десериализации ответа")
	}

	matchString, err := regexp.MatchString("^test/.+$", resp.Url)
	if err != nil || !matchString {
		t.Errorf("url неверен: получили %s, а хотели %s", rr.Body, "test/?")
	}
}
