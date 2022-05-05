package handler

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/onorbit/pixelite/library"
)

func mountLibrary(c echo.Context) error {
	rootPath := c.FormValue("rootPath")
	if len(rootPath) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := library.MountLibrary(rootPath); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	// TODO : the API should return appropriate response.
	return c.NoContent(http.StatusOK)
}

func getLibrary(c echo.Context) error {
	id := c.Param("id")
	library := library.GetLibrary(id)
	if library == nil {
		return c.NoContent(http.StatusNotFound)
	}

	libDesc := library.Describe()
	return c.JSON(http.StatusOK, libDesc)
}

func unmountLibrary(c echo.Context) error {
	id := c.Param("id")
	if err := library.UnmountLibrary(id); err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.NoContent(http.StatusOK)
}

func rescanLibrary(c echo.Context) error {
	id := c.Param("id")
	library := library.GetLibrary(id)
	if library == nil {
		return c.NoContent(http.StatusNotFound)
	}

	if err := library.Rescan(); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func setLibraryTitle(c echo.Context) error {
	id := c.Param("id")
	library := library.GetLibrary(id)
	if library == nil {
		return c.NoContent(http.StatusNotFound)
	}

	title := c.FormValue("title")
	if len(title) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := library.SetTitle(title); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func listLibrary(c echo.Context) error {
	list := library.ListLibrary()
	return c.JSON(http.StatusOK, list)
}
