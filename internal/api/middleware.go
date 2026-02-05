package api

import (
	"net/http"
	"slices"
	"strings"
)

type MiddlewareChain []func(http.Handler) http.Handler

func (c MiddlewareChain) thenFunc(h http.HandlerFunc) http.Handler {
	return c.then(h)
}

func (c MiddlewareChain) then(h http.Handler) http.Handler {
	for _, mdw := range slices.Backward(c) {
		h = mdw(h)
	}
	return h
}

var middleware = MiddlewareChain{
	recoverer,
	realIP,
	helmet,
}

func recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func realIP(next http.Handler) http.Handler {
	var (
		trueClientIP  = http.CanonicalHeaderKey("True-Client-IP")
		xRealIP       = http.CanonicalHeaderKey("X-Real-IP")
		xForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
	)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ip string

		if tcip := r.Header.Get(trueClientIP); tcip != "" {
			ip = tcip
		} else if xrip := r.Header.Get(xRealIP); xrip != "" {
			ip = xrip
		} else if xff := r.Header.Get(xForwardedFor); xff != "" {
			i := strings.Index(xff, ", ")
			if i == -1 {
				i = len(xff)
			}
			ip = xff[:i] // first IP in the comma-separated list
		}

		if ip != "" {
			r.RemoteAddr = ip
		}

		next.ServeHTTP(w, r)
	})
}

func helmet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent sensitive information from being cached.
		w.Header().Set("Cache-Control", "no-store")

		// To protect against drag-and-drop style clickjacking attacks.
		w.Header().Set("Content-Security-Policy", "frame-ancestors 'none'")
		w.Header().Set("X-Frame-Options", "DENY")

		// To prevent browsers from performing MIME sniffing, and inappropriately interpreting responses as HTML.
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Require connections over HTTPS and to protect against spoofed certificates.
		w.Header().Set("Strict-Transport-Security", "max-age=31536000")

		next.ServeHTTP(w, r)
	})
}
