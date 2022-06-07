package handler

import (
	"net/http"

	"github.com/labstack/echo"
)

type DirectoryEntryType byte

const (
	Directory = iota
	ImageFile
)

type DirectoryEntry struct {
	Name string             `json:"name"`
	Type DirectoryEntryType `json:"type"`
}

func handleIndex(c echo.Context) error {
	return c.Redirect(http.StatusSeeOther, "/libraries")
}

func Initialize(listenPort int) error {
	if err := initRouter(listenPort); err != nil {
		return err
	}

	return nil
}

func Cleanup() error {
	if err := cleanupRouter(); err != nil {
		return err
	}

	return nil
}
