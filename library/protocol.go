package library

type LibraryDesc struct {
	Id     string   `json:"id"`
	Desc   string   `json:"desc"`
	Albums []string `json:"albums"`
}

type LibrarySummeryDesc struct {
	Id   string `json:"id"`
	Desc string `json:"desc"`
}
