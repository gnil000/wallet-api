package utils

import (
	"net/http"
	"strconv"

	"github.com/goccy/go-json"
)

func RespondStatus(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(status)
	_, _ = w.Write(response)
}

func RespondError(w http.ResponseWriter, code int, message string) {
	RespondJSON(w, code, map[string]string{"code": strconv.Itoa(code), "error": message})
}
