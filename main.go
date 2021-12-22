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
	e.GET("/apis/list/:libid/:albumid", handler.ListPath)
	e.GET("/apis/thumbnail/:libid/:albumid/:filename", handler.ServeThumbnail)
	e.GET("/apis/image/:libid/:albumid/:filename", handler.ServeImage)

	e.POST("/apis/library", handler.CreateLibrary)
	e.GET("/apis/library/:id", handler.GetLibrary)
	e.DELETE("/apis/library/:id", handler.DeleteLibrary)
	e.GET("/apis/libraries", handler.ListLibrary)

	e.Static("/statics", "statics")
	e.File("/libraries", "views/libraries.html")
	e.File("/library/:id", "views/library.html")
	e.File("/thumbnails/:libid/:albumid", "views/thumbnails.html")

	e.Logger.Fatal(e.Start(":10900"))
}
