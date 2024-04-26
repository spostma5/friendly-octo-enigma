package api

import (
	"log/slog"
	"net/http"
)

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Handling request", "Method", r.Method, "Path", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// Where to put Auth stuff if needed in future
