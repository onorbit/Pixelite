package handler

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/labstack/echo"
	"github.com/onorbit/pixelite/image"
	"github.com/onorbit/pixelite/library"
	"github.com/onorbit/pixelite/thumbnail"
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
	// acquire Library and belonging Album from input.
	libraryID := c.Param("libid")
	albumID := c.Param("albumid")

	albumID, err := url.QueryUnescape(albumID)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	targetLibrary := library.GetLibrary(libraryID)
	if targetLibrary == nil {
		return c.NoContent(http.StatusNotFound)
	}

	targetAlbum := targetLibrary.GetAlbum(albumID)
	if targetAlbum == nil {
		return c.NoContent(http.StatusNotFound)
	}

	// list content of album path.
	albumPath := targetAlbum.GetPath()
	content, err := ioutil.ReadDir(albumPath)
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
	// prepare parameters.
	libraryID := c.Param("libid")

	albumID := c.Param("albumid")
	albumID, err := url.QueryUnescape(albumID)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	fileName := c.Param("filename")
	fileName, err = url.QueryUnescape(fileName)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	// compose target file path.
	targetLibrary := library.GetLibrary(libraryID)
	if targetLibrary == nil {
		return c.NoContent(http.StatusNotFound)
	}

	targetAlbum := targetLibrary.GetAlbum(albumID)
	if targetAlbum == nil {
		return c.NoContent(http.StatusNotFound)
	}

	filePath := filepath.Join(targetAlbum.GetPath(), fileName)
	thumbnailPath := thumbnail.GetThumbnailPath(filePath)
	if len(thumbnailPath) == 0 {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.File(thumbnailPath)
}

func ServeImage(c echo.Context) error {
	// prepare parameters.
	libraryID := c.Param("libid")

	albumID := c.Param("albumid")
	albumID, err := url.QueryUnescape(albumID)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	fileName := c.Param("filename")
	fileName, err = url.QueryUnescape(fileName)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	// compose target file path.
	targetLibrary := library.GetLibrary(libraryID)
	if targetLibrary == nil {
		return c.NoContent(http.StatusNotFound)
	}

	targetAlbum := targetLibrary.GetAlbum(albumID)
	if targetAlbum == nil {
		return c.NoContent(http.StatusNotFound)
	}

	filePath := filepath.Join(targetAlbum.GetPath(), fileName)
	return c.File(filePath)
}
