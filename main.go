package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"azurestud.io/pixelite/config"
	"azurestud.io/pixelite/thumbnail"
	"github.com/labstack/echo"
)

func main() {
	if err := config.Initialize("pixelite.json"); err != nil {
		panic(err)
	}
	if err := thumbnail.Initialize(); err != nil {
		panic(err)
	}

	e := echo.New()
	e.GET("/apis/list/*", listPath)
	e.GET("/apis/thumbnail/*", serveThumbnail)

	e.Static("/statics", "statics")
	e.File("/thumbnails/*", "views/thumbnails.html")

	e.Logger.Fatal(e.Start(":10900"))
}

type DirectoryEntryType byte

const (
	Directory = iota
	ImageFile
)

type DirectoryEntry struct {
	Name string             `json:"name"`
	Type DirectoryEntryType `json:"type"`
}

var imageExt = []string{
	".jpg",
	".png",
}

func isImageFile(fileName string) bool {
	fileExt := filepath.Ext(fileName)
	fileExt = strings.ToLower(fileExt)
	for _, ext := range imageExt {
		if ext == fileExt {
			return true
		}
	}

	return false
}

func listPath(c echo.Context) error {
	subPath := c.Param("*")
	if len(subPath) <= 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	subPath, err := url.PathUnescape(subPath)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	rootPath := config.Get().RootPath
	targetPath := filepath.Join(rootPath, subPath)
	content, err := ioutil.ReadDir(targetPath)
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
		} else if isImageFile(entry.Name()) {
			newEntry := DirectoryEntry{
				Name: entry.Name(),
				Type: ImageFile,
			}
			entryList = append(entryList, newEntry)
		}
	}

	return c.JSON(http.StatusOK, entryList)
}

func serveThumbnail(c echo.Context) error {
	subPath := c.Param("*")
	if len(subPath) <= 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	subPath, err := url.PathUnescape(subPath)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	rootPath := config.Get().RootPath
	imgPath := filepath.Join(rootPath, subPath)
	thumbnailPath := thumbnail.GetThumbnailPath(imgPath)
	if len(thumbnailPath) == 0 {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.File(thumbnailPath)
}
