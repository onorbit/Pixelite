package handler

import (
	"net/http"

	"azurestud.io/pixelite/library"
	"github.com/labstack/echo"
)

func CreateLibrary(c echo.Context) error {
	rootPath := c.FormValue("rootPath")
	if len(rootPath) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := library.CreateLibrary(rootPath); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func ListLibrary(c echo.Context) error {
	list := library.ListLibrary()
	return c.JSON(http.StatusOK, list)
}

func GetLibrary(c echo.Context) error {
	id := c.Param("id")
	library := library.GetLibrary(id)
	if library == nil {
		return c.NoContent(http.StatusNotFound)
	}

	libDesc := library.Describe()
	return c.JSON(http.StatusOK, libDesc)
}
