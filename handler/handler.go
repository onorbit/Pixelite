package handler

import (
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/labstack/echo"
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

func handleIndex(c echo.Context) error {
	return c.Redirect(http.StatusSeeOther, "/libraries")
}

func handleListPath(c echo.Context) error {
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
	imageList, err := targetAlbum.ListImages()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	entryList := make([]DirectoryEntry, 0, len(imageList))
	for _, fileName := range imageList {
		newEntry := DirectoryEntry{
			Name: fileName,
			Type: ImageFile,
		}

		entryList = append(entryList, newEntry)
	}

	return c.JSON(http.StatusOK, entryList)
}

func handleServeThumbnail(c echo.Context) error {
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

	thumbnailPath := thumbnail.GetThumbnailPath(fileName, targetAlbum.GetPath(), albumID, libraryID)
	if len(thumbnailPath) == 0 {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.File(thumbnailPath)
}

func handleServeImage(c echo.Context) error {
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
