package handler

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"azurestud.io/pixelite/image"
	"azurestud.io/pixelite/thumbnail"
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

func ListPath(c echo.Context) error {
	subPath := c.Param("*")
	if len(subPath) <= 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	subPath, err := url.PathUnescape(subPath)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	content, err := ioutil.ReadDir(subPath)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	entryList := make([]DirectoryEntry, 0, len(content))
	for _, entry := range content {
		if entry.IsDir() {
			newEntry := DirectoryEntry{
				Name: entry.Name(),
				Type: Directory,
			}
			entryList = append(entryList, newEntry)
		} else if image.IsImageFile(entry.Name()) {
			newEntry := DirectoryEntry{
				Name: entry.Name(),
				Type: ImageFile,
			}
			entryList = append(entryList, newEntry)
		}
	}

	return c.JSON(http.StatusOK, entryList)
}

func ServeThumbnail(c echo.Context) error {
	subPath := c.Param("*")
	if len(subPath) <= 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	subPath, err := url.PathUnescape(subPath)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	thumbnailPath := thumbnail.GetThumbnailPath(subPath)
	if len(thumbnailPath) == 0 {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.File(thumbnailPath)
}
