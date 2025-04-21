package api

import (
	"encoding/json"
	"net/http"
)

const (
	StatusOk  = "OK"
	StatusErr = "ERROR"
)

// Response базовая структура все api ответов
type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func Ok() Response {
	return Response{Status: StatusOk}
}

func Err(err string) Response {
	return Response{Status: StatusErr, Error: err}
}

// WriteRespJson записывает результат в переданный ResponseWriter в формате json
func WriteRespJson(w http.ResponseWriter, data any, code int) {
	d, _ := json.Marshal(data)
	w.WriteHeader(code)
	write, err := w.Write(d)

	if err != nil {
		WriteRespJson(w, Err(err.Error()), http.StatusInternalServerError)
		return
	}

	_ = write
}

// ResponseOk формирует ответ и записывает результат в переданный http.ResponseWriter
func ResponseOk(w http.ResponseWriter, url string, status int) Resp {
	r := Resp{
		Response: Ok(),
		Url:      url,
	}

	WriteRespJson(w, r, status)

	return r
}
