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
	e.HideBanner = true

	// general APIs.
	e.GET("/apis/list/:libid/:albumid", handleListPath)
	e.GET("/apis/thumbnail/:libid/:albumid/:filename", handleServeThumbnail)
	e.GET("/apis/image/:libid/:albumid/:filename", handleServeImage)

	// library APIs.
	e.POST("/apis/library", mountLibrary)
	e.GET("/apis/library/:id", getLibrary)
	e.DELETE("/apis/library/:id", unmountLibrary)
	e.POST("/apis/library/:id/rescan", rescanLibrary)
	e.POST("/apis/library/:id/title", setLibraryTitle)
	e.GET("/apis/libraries", listLibrary)

	// frontend routes.
	e.GET("/", handleIndex)
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
