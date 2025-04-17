package url

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"
	"urlShortener/internal/config"
	"urlShortener/internal/storage/memdb"
)

var a *API

func TestMain(m *testing.M) {
	cfg := config.Config{
		Env: "test",
		Http: config.HttpServer{
			Address:     "test",
			Timeout:     0,
			IdleTimeout: 0,
		},
	}

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	a = New(memdb.New(&cfg), &cfg.Http, log)
	a.Router = mux.NewRouter()
	a.Endpoints()
	os.Exit(m.Run())
}

func TestAPI_Save(t *testing.T) {
	data := SaveUrlReq{
		Url: "http://google.com",
	}
	payload, _ := json.Marshal(data)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/url", bytes.NewBuffer(payload))

	rr := httptest.NewRecorder()

	a.Router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}

	var resp UrlResp
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Errorf("ошибка при десериализации ответа")
	}

	matchString, err := regexp.MatchString("^test/.+$", resp.Url)
	if err != nil || !matchString {
		t.Errorf("url неверен: получили %s, а хотели %s", rr.Body, "test/?")
	}
}

func TestAPI_Original(t *testing.T) {
	const url = "http://yandex.ru"
	const shortUrl = "test/2"

	_, err := a.db.SaveUrl(url, shortUrl)
	if err != nil {
		t.Errorf("ошибка при подготовке данных")
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/url/original", nil)
	q := req.URL.Query()
	q.Add("url", shortUrl)
	req.URL.RawQuery = q.Encode()

	rr := httptest.NewRecorder()

	a.Router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}

	var resp UrlResp
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Errorf("ошибка при десериализации ответа")
	}

	if url != resp.Url {
		t.Errorf("url неверен: получили %s, а хотели %s", resp.Url, url)
	}
}
