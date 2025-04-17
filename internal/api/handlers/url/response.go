package url

import (
	"net/http"
	"urlShortener/internal/api/handlers"
)

// Resp структура ответа для хендлеров url
type Resp struct {
	handlers.Response
	Url string `json:"url"`
}

// ResponseOk формирует ответ и записывает результат в переданный http.ResponseWriter
func ResponseOk(w http.ResponseWriter, url string, status int) {
	r := Resp{
		Response: handlers.Ok(),
		Url:      url,
	}

	handlers.WriteRespJson(w, r, status)
}
