package thumbnail

import (
	"os"
	"sync"

	"github.com/onorbit/pixelite/config"
)

type thumbnailedAlbumKey struct {
	libraryID string
	albumID   string
}

func Initialize() error {
	cfg := config.Get()

	// make paths.
	if err := os.MkdirAll(cfg.Thumbnail.StorePath, 0700); err != nil {
		return err
	}
	if err := os.MkdirAll(cfg.Cover.StorePath, 0700); err != nil {
		return err
	}

	// initialize album manager.
	gAlbumManager = albumManager{
		thumbnailedAlbums:  make(map[thumbnailedAlbumKey]*thumbnailedAlbum),
		recentAccessAlbums: make(map[thumbnailedAlbumKey]struct{}),
		mutex:              sync.Mutex{},
	}

	if err := gAlbumManager.loadFromGlobalDB(); err != nil {
		return err
	}

	gAlbumManager.startTick()

	// initialze cover manager.
	gCoverManager = coverManager{
		covers: make(map[thumbnailedAlbumKey]string),
	}

	if err := gCoverManager.loadFromGlobalDB(); err != nil {
		return err
	}

	return nil
}

func Cleanup() {
	gAlbumManager.stopTick()
}

func GetThumbnailPath(fileName, albumPath, albumID, libraryID string) string {
	return gAlbumManager.getThumbnailPath(fileName, albumPath, albumID, libraryID)
}

func GetAlbumCover(fileName, albumPath, albumID, libraryID string) string {
	return gCoverManager.getAlbumCover(fileName, albumPath, albumID, libraryID)
}
