package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)

func JSONEncode(w http.ResponseWriter, v any, status int) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		return errors.New("Unable to encode JSON")
	}

	return nil
}
