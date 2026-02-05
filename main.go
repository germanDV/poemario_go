package main

import (
	"errors"
	"net/http"

	"github.com/germandv/poemario/internal/api"
	"github.com/germandv/poemario/internal/poemario"
)

func main() {
	err := api.New(poemario.NewClient(), 8080).ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}
