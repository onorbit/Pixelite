package main

import (
	"math/rand"
	"time"

	"github.com/labstack/echo"
	"github.com/onorbit/pixelite/config"
	"github.com/onorbit/pixelite/globaldb"
	"github.com/onorbit/pixelite/handler"
	"github.com/onorbit/pixelite/library"
	"github.com/onorbit/pixelite/thumbnail"
)

func main() {
	rand.Seed(time.Now().UnixNano())

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

	e.POST("/apis/library", handler.CreateLibrary)
	e.GET("/apis/library/:id", handler.GetLibrary)
	e.DELETE("/apis/library/:id", handler.DeleteLibrary)
	e.GET("/apis/libraries", handler.ListLibrary)

	e.Static("/statics", "statics")
	e.File("/libraries", "views/libraries.html")
	e.File("/library/:id", "views/library.html")
	e.File("/thumbnails/*", "views/thumbnails.html")

	e.Logger.Fatal(e.Start(":10900"))
}
