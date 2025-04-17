package api

import (
	"encoding/json"
	"net/http"
)

const (
	StatusOk  = "OK"
	StatusErr = "ERROR"
)

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

func WriteResp(w http.ResponseWriter, data any, code int) {
	d, _ := json.Marshal(data)
	w.WriteHeader(code)
	write, err := w.Write(d)

	if err != nil {
		WriteResp(w, Err(err.Error()), http.StatusInternalServerError)
		return
	}

	_ = write
}
