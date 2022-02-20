package handler

import (
	"context"
	"sync"

	"github.com/labstack/echo"
)

var gCancelFunc context.CancelFunc
var gWaitGroup sync.WaitGroup
var gEcho *echo.Echo

func initRouter() error {
	e := echo.New()

	e.GET("/apis/list/:libid/:albumid", listPath)
	e.GET("/apis/thumbnail/:libid/:albumid/:filename", serveThumbnail)
	e.GET("/apis/image/:libid/:albumid/:filename", serveImage)

	e.POST("/apis/library", createLibrary)
	e.GET("/apis/library/:id", getLibrary)
	e.DELETE("/apis/library/:id", deleteLibrary)
	e.GET("/apis/libraries", listLibrary)

	e.Static("/statics", "frontend/statics")
	e.File("/libraries", "frontend/views/libraries.html")
	e.File("/library/:id", "frontend/views/library.html")
	e.File("/thumbnails/:libid/:albumid", "frontend/views/thumbnails.html")

	gEcho = e

	gWaitGroup.Add(1)
	go runFunc()

	return nil
}

func runFunc() {
	// execution blocks here.
	gEcho.Start(":10900")

	// following code is executed on shutdown.
	gWaitGroup.Done()
}

func cleanupRouter() error {
	gEcho.Shutdown(context.Background())
	gWaitGroup.Wait()

	return nil
}
