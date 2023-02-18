package server

import (
	"encoding/json"
	"net/http"
)

type response[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func responseOk[T any](w http.ResponseWriter, data T, message ...string) {
	r := response[T]{
		Message: "OK",
		Data:    data,
	}

	if len(message) == 1 {
		r.Message = message[0]
	}

	raw, err := json.Marshal(r)
	if err != nil {
		responseErr(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(raw); err != nil {
		panic(err)
	}
}

func responseErr(w http.ResponseWriter, err error, status ...int) {
	r := response[error]{
		Message: err.Error(),
		Data:    nil,
	}

	raw, _ := json.Marshal(r)
	s := http.StatusInternalServerError
	if len(status) == 1 {
		s = status[0]
	}

	w.WriteHeader(s)
	if _, err := w.Write(raw); err != nil {
		panic(err)
	}
}
