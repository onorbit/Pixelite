package main

import (
	"azurestud.io/pixelite/config"
	"azurestud.io/pixelite/globaldb"
	"azurestud.io/pixelite/handler"
	"azurestud.io/pixelite/library"
	"azurestud.io/pixelite/thumbnail"
	"github.com/labstack/echo"
)

func main() {
	if err := config.Initialize("pixelite.json"); err != nil {
		panic(err)
	}
	if err := globaldb.Initialize(config.Get().GlobalDBPath); err != nil {
		panic(err)
	}
	if err := thumbnail.Initialize(); err != nil {
		panic(err)
	}
	if err := library.Initialize(); err != nil {
		panic(err)
	}

	e := echo.New()
	e.GET("/apis/list/*", handler.ListPath)
	e.GET("/apis/thumbnail/*", handler.ServeThumbnail)
	e.GET("/apis/libraries", handler.ListLibrary)
	e.GET("/apis/library/:id", handler.GetLibrary)
	e.POST("/apis/library", handler.CreateLibrary)

	e.Static("/statics", "statics")
	e.File("/thumbnails/*", "views/thumbnails.html")

	e.Logger.Fatal(e.Start(":10900"))
}
