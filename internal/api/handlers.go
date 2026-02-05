package api

import (
	"encoding/json"
	"net/http"

	"github.com/germandv/poemario/internal/errors"
	"github.com/germandv/poemario/internal/poemario"
)

func (api API) GetAuthors(w http.ResponseWriter, r *http.Request) error {
	type GetAuthorsResp struct {
		Authors []poemario.Author `json:"authors"`
		Total   int               `json:"total"`
	}

	authors, err := api.Poemario.GetAuthors()
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(GetAuthorsResp{
		Authors: authors,
		Total:   len(authors),
	})
}

func (api API) GetPoems(w http.ResponseWriter, r *http.Request) error {
	type GetPoemsResp struct {
		Author poemario.Author `json:"author"`
		Titles []string        `json:"titles"`
		Total  int             `json:"total"`
	}

	poems, err := api.Poemario.GetPoems(r.PathValue("author_name"))
	if err != nil {
		return err
	}

	if len(poems) == 0 {
		return errors.NewNotFoundError("no poems found for author")
	}

	titles := make([]string, 0, len(poems))
	for _, poem := range poems {
		titles = append(titles, poem.Title)
	}

	return json.NewEncoder(w).Encode(GetPoemsResp{
		Author: poems[0].Author,
		Titles: titles,
		Total:  len(poems),
	})
}

func (api API) GetPoem(w http.ResponseWriter, r *http.Request) error {
	poem, err := api.Poemario.GetPoem(r.PathValue("author_name"), r.PathValue("title"))
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(poem)
}

func (api API) GetRandomPoem(w http.ResponseWriter, r *http.Request) error {
	poem, err := api.Poemario.GetRandomPoem()
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(poem)
}
