package api

import (
	"net/http"
	"time"
)

type IP string

type RequestCount struct {
	count      int
	lastUpdate time.Time
}

func rateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ip := r.RemoteAddr

		next.ServeHTTP(w, r)
	})
}
