package poemario

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/germandv/poemario/internal/errors"
)

type Client struct {
	baseURL    string
	httpClient http.Client
}

type apiErrorResponse struct {
	Status int `json:"status"`
}

func checkResponseForError(data []byte) error {
	var errResp apiErrorResponse

	err := json.Unmarshal(data, &errResp)
	if err != nil {
		// We could not decode the response into an error response,
		// so we assume it is not an error
		return nil
	}

	if errResp.Status >= 400 {
		return fmt.Errorf("received status %d", errResp.Status)
	}

	return nil
}

func NewClient() Client {
	return Client{
		baseURL:    "https://poetrydb.org",
		httpClient: http.Client{Timeout: 5 * time.Second},
	}
}

func (c Client) GetAuthors() ([]Author, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/author", c.baseURL),
		nil,
	)
	if err != nil {
		return nil, errors.NewInternalServerError(err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.NewServiceUnavailableError("failed to reach poetry service")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.NewServiceUnavailableError(fmt.Sprintf("poetry service returned status %s", resp.Status))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewInternalServerError(err)
	}
	if err := checkResponseForError(data); err != nil {
		return nil, errors.NewServiceUnavailableError("poetry service error")
	}

	body := map[string][]string{}
	err = json.Unmarshal(data, &body)
	if err != nil {
		return nil, errors.NewInternalServerError(err)
	}

	authors := make([]Author, 0, len(body["authors"]))
	for _, author := range body["authors"] {
		authors = append(authors, Author{FullName: author})
	}

	return authors, nil
}

func (c Client) GetPoems(author string) ([]Poem, error) {
	if author == "" || len(author) > 100 {
		return nil, errors.NewValidationError("author name must be between 1 and 100 characters")
	}

	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/author/%s/title,author", c.baseURL, url.PathEscape(author)),
		nil,
	)
	if err != nil {
		return nil, errors.NewInternalServerError(err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.NewServiceUnavailableError("failed to reach poetry service")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.NewServiceUnavailableError(fmt.Sprintf("poetry service returned status %s", resp.Status))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewInternalServerError(err)
	}
	if err := checkResponseForError(data); err != nil {
		return nil, errors.NewServiceUnavailableError("poetry service error")
	}

	body := []struct {
		Title  string `json:"title"`
		Author string `json:"author"`
	}{}

	err = json.Unmarshal(data, &body)
	if err != nil {
		return nil, errors.NewInternalServerError(err)
	}

	var matchedAuthor Author
	poems := make([]Poem, 0, len(body))

	for i, entry := range body {
		if i == 0 {
			matchedAuthor = Author{FullName: entry.Author}
		} else if entry.Author != matchedAuthor.FullName {
			return nil, errors.NewConflictError(fmt.Sprintf("matched more than one author for search: %q, try something more specific", author))
		}
		poems = append(poems, Poem{Author: matchedAuthor, Title: entry.Title})
	}

	return poems, nil
}

func (c Client) GetPoem(author string, title string) (Poem, error) {
	if author == "" || len(author) > 100 {
		return Poem{}, errors.NewValidationError("author name must be between 1 and 100 characters")
	}
	if title == "" || len(title) > 200 {
		return Poem{}, errors.NewValidationError("title must be between 1 and 200 characters")
	}

	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/author,title/%s;%s", c.baseURL, url.PathEscape(author), url.PathEscape(title)),
		nil,
	)
	if err != nil {
		return Poem{}, errors.NewInternalServerError(err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Poem{}, errors.NewServiceUnavailableError("failed to reach poetry service")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Poem{}, errors.NewServiceUnavailableError(fmt.Sprintf("poetry service returned status %s", resp.Status))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Poem{}, errors.NewInternalServerError(err)
	}
	if err := checkResponseForError(data); err != nil {
		return Poem{}, errors.NewServiceUnavailableError("poetry service error")
	}

	body := []struct {
		Author string   `json:"author"`
		Title  string   `json:"title"`
		Lines  []string `json:"lines"`
	}{}

	err = json.Unmarshal(data, &body)
	if err != nil {
		return Poem{}, errors.NewInternalServerError(err)
	}

	if len(body) == 0 {
		return Poem{}, errors.NewNotFoundError("poem not found")
	}

	return Poem{
		Author: Author{FullName: body[0].Author},
		Title:  body[0].Title,
		Lines:  body[0].Lines,
	}, nil
}

func (c Client) GetRandomPoem() (Poem, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/random", c.baseURL),
		nil,
	)
	if err != nil {
		return Poem{}, errors.NewInternalServerError(err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Poem{}, errors.NewServiceUnavailableError("failed to reach poetry service")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Poem{}, errors.NewServiceUnavailableError(fmt.Sprintf("poetry service returned status %s", resp.Status))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Poem{}, errors.NewInternalServerError(err)
	}
	if err := checkResponseForError(data); err != nil {
		return Poem{}, errors.NewServiceUnavailableError("poetry service error")
	}

	body := []struct {
		Author string   `json:"author"`
		Title  string   `json:"title"`
		Lines  []string `json:"lines"`
	}{}

	err = json.Unmarshal(data, &body)
	if err != nil {
		return Poem{}, errors.NewInternalServerError(err)
	}

	if len(body) == 0 {
		return Poem{}, errors.NewNotFoundError("no random poems found")
	}

	return Poem{
		Author: Author{FullName: body[0].Author},
		Title:  body[0].Title,
		Lines:  body[0].Lines,
	}, nil
}
