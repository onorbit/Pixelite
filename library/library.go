package library

type Library struct {
	RootPath string
	Albums   map[string]struct{}
}
