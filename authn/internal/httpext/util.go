package httpext

import (
	"encoding/json"
	"net/http"
)

func JSONError(w http.ResponseWriter, err interface{}, code int) {
	h := w.Header()
	h.Del("Content-Length")
	h.Set("Content-Type", "application/json; charset=utf-8")
	h.Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(err)
}
