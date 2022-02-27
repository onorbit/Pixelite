package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/onorbit/pixelite/config"
	"github.com/onorbit/pixelite/database/globaldb"
	"github.com/onorbit/pixelite/handler"
	"github.com/onorbit/pixelite/library"
	"github.com/onorbit/pixelite/thumbnail"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	conf := config.Get()

	if err := config.Initialize("pixelite.json"); err != nil {
		panic(err)
	}
	if err := globaldb.Initialize(conf.GlobalDBPath); err != nil {
		panic(err)
	}
	if err := thumbnail.Initialize(); err != nil {
		panic(err)
	}
	if err := library.Initialize(); err != nil {
		panic(err)
	}
	if err := handler.Initialize(conf.ListenPort); err != nil {
		panic(err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	fmt.Printf("shutting down.\n")

	handler.Cleanup()
}
