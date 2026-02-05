package poemario

type Poem struct {
	Author Author   `json:"author"`
	Title  string   `json:"title"`
	Lines  []string `json:"lines"`
}
