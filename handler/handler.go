package handler

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/onorbit/pixelite/config"
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

type Configs struct {
	ThumbnailSize int `json:"thumbnailSize"`
	CoverSize     int `json:"coverSize"`
}

func handleIndex(c echo.Context) error {
	return c.Redirect(http.StatusSeeOther, "/libraries")
}

func getConfigs(c echo.Context) error {
	conf := config.Get()

	ret := Configs{
		ThumbnailSize: conf.Thumbnail.MaxDimension,
		CoverSize:     conf.Cover.MaxDimension,
	}
	return c.JSON(http.StatusOK, ret)
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
