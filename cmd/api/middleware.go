package api

import (
	"log/slog"
	"net/http"
)

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Handling request", "Method", r.Method, "Path", r.URL.Path)
		next.ServeHTTP(w, r)

		// Thought about logging here if we have erroneous data or errors here,
		// but given I'm not sure really what is going to be stored here
		// and don't want to be logging anything sensitive I'll leave
		// it for now
	})
}

// Where to put Auth stuff if needed in future, although if we add much more to this
// I might just put it in a separate package.
// Would also need to setup some middleware chaining in that case
