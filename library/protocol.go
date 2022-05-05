package library

type LibraryDesc struct {
	Id     string   `json:"id"`
	Title  string   `json:"title"`
	Albums []string `json:"albums"`
}

type LibrarySummeryDesc struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}
