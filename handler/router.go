package handler

import (
	"context"
	"fmt"
	"sync"

	"github.com/labstack/echo"
)

var gCancelFunc context.CancelFunc
var gWaitGroup sync.WaitGroup
var gEcho *echo.Echo

func initRouter(listenPort int) error {
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
	go runFunc(listenPort)

	return nil
}

func runFunc(listenPort int) {
	// execution blocks here.
	addr := fmt.Sprintf(":%d", listenPort)
	gEcho.Start(addr)

	// following code is executed on shutdown.
	gWaitGroup.Done()
}

func cleanupRouter() error {
	gEcho.Shutdown(context.Background())
	gWaitGroup.Wait()

	return nil
}
