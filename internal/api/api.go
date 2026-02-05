package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/germandv/poemario/internal/errors"
	"github.com/germandv/poemario/internal/poemario"
)

type API struct {
	Poemario poemario.Client
	Port     int
}

func New(poemario poemario.Client, port int) API {
	return API{
		Poemario: poemario,
		Port:     port,
	}
}

func (api API) ListenAndServe() error {
	mux := &http.ServeMux{}
	mux.Handle("GET /health", middleware.thenFunc(api.Health))
	mux.Handle("GET /v1/authors", middleware.thenFunc(api.errorHandler(api.GetAuthors)))
	mux.Handle("GET /v1/authors/{author_name}/poems", middleware.thenFunc(api.errorHandler(api.GetPoems)))
	mux.Handle("GET /v1/authors/{author_name}/poems/{title}", middleware.thenFunc(api.errorHandler(api.GetPoem)))
	mux.Handle("GET /v1/poems/random", middleware.thenFunc(api.errorHandler(api.GetRandomPoem)))

	fmt.Printf("Starting server on :%d\n", api.Port)
	fmt.Println("Routes:")
	fmt.Println("  - GET /health")
	fmt.Println("  - GET /v1/authors")
	fmt.Println("  - GET /v1/authors/{author_name}/poems")
	fmt.Println("  - GET /v1/authors/{author_name}/poems/{title}")
	fmt.Println("  - GET /v1/poems/random")

	server := http.Server{
		Addr:              fmt.Sprintf(":%d", api.Port),
		Handler:           mux,
		IdleTimeout:       1 * time.Minute,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
	}

	return server.ListenAndServe()
}

// errorHandler wraps an HTTP handler that returns an error and maps it to HTTP responses
func (api API) errorHandler(h func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := h(w, r)
		if err != nil {
			if appErr, ok := err.(*errors.AppError); ok {
				w.WriteHeader(appErr.Code)
				json.NewEncoder(w).Encode(map[string]string{"error": appErr.Message})
				return
			}

			fmt.Printf("Unexpected error: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		}
	}
}

func (api API) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
