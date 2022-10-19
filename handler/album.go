package handler

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo"
	"github.com/onorbit/pixelite/library"
	"github.com/onorbit/pixelite/media"
	"github.com/onorbit/pixelite/thumbnail"
)

func getAlbumImageList(c echo.Context) error {
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
	imageList, err := targetAlbum.ListMedias()
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

func getAlbumImage(c echo.Context) error {
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

	// check if the fileName is supported media file type.
	if !media.IsSupportedMedia(fileName) {
		return c.NoContent(http.StatusBadRequest)
	}

	// check if the fileName contains PathSeparator, which may lead to other path.
	if strings.ContainsRune(fileName, os.PathSeparator) {
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

func getAlbumCover(c echo.Context) error {
	// prepare parameters.
	libraryID := c.Param("libid")

	albumID := c.Param("albumid")
	albumID, err := url.QueryUnescape(albumID)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	// get cover image path.
	targetLibrary := library.GetLibrary(libraryID)
	if targetLibrary == nil {
		return c.NoContent(http.StatusNotFound)
	}

	targetAlbum := targetLibrary.GetAlbum(albumID)
	if targetAlbum == nil {
		return c.NoContent(http.StatusNotFound)
	}

	coverPath := thumbnail.GetAlbumCover(targetAlbum.GetCoverFileName(), targetAlbum.GetPath(), albumID, libraryID)
	if len(coverPath) == 0 {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.File(coverPath)
}

func getAlbumImageThumbnail(c echo.Context) error {
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
