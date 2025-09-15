package errh

import (
	"encoding/json"
	"net/http"
)

type typedHTTPError struct {
	Error string `json:"error"`
	Type  string `json:"error_type,omitempty"`
}

func setHTTPError(w http.ResponseWriter, code int, msg, errType string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Del("Cache-Control")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(&typedHTTPError{
		Error: msg,
		Type:  errType,
	})
	if err != nil {
		w.WriteHeader(defaultHTTPCode)
	}
}
