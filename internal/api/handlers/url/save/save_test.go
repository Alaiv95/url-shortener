package save_test

//
//import (
//	"bytes"
//	"context"
//	"encoding/json"
//	"log/slog"
//	"net/http"
//	"net/http/httptest"
//	"os"
//	"strings"
//	"testing"
//	"time"
//	"urlShortener/internal/api/handlers/url"
//	"urlShortener/internal/api/handlers/url/save"
//	"urlShortener/internal/config"
//	"urlShortener/internal/storage/memdb"
//)
//
//var log *slog.Logger
//var saver save.UrlSaver
//var cfg config.Config
//
//const domain = "https://test.ru"
//
//type DummyProducer struct{}
//type DummySetter struct{}
//
//func (p DummyProducer) Produce(_ []byte, _ context.Context) {}
//func (s DummySetter) Set(_ string, _ any, _ time.Duration) error {
//	return nil
//}
//
//func TestMain(m *testing.M) {
//	cfg = config.Config{
//		Env: "test",
//		Http: config.HttpServer{
//			Address:     "test",
//			Timeout:     0,
//			IdleTimeout: 0,
//		},
//	}
//
//	log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
//		Level: slog.LevelDebug,
//	}))
//
//	saver = memdb.New(&cfg)
//
//	os.Exit(m.Run())
//}
//
//func TestAPI_Save(t *testing.T) {
//	data := save.UrlReq{
//		Url: "http://google.com",
//	}
//	payload, _ := json.Marshal(data)
//	h := save.New(log, saver, DummyProducer{}, domain, DummySetter{})
//
//	req := httptest.NewRequest(http.MethodPost, "/api/v1/url", bytes.NewBuffer(payload))
//
//	rr := httptest.NewRecorder()
//
//	h.ServeHTTP(rr, req)
//	if rr.Code != http.StatusOK {
//		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
//	}
//
//	var resp url.Resp
//	err := json.Unmarshal(rr.Body.Bytes(), &resp)
//	if err != nil {
//		t.Errorf("ошибка при десериализации ответа")
//	}
//
//	if err != nil || resp.Url == "" {
//		t.Errorf("url неверен: получили %s, а хотели %s", rr.Body, "?")
//	}
//
//	if !strings.Contains(resp.Url, domain) {
//		t.Errorf("url неверен: получили %s, а хотели %s", rr.Body, domain+"/?")
//	}
//}
