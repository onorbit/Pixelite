package thumbnail

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/onorbit/pixelite/config"
	"github.com/onorbit/pixelite/database/globaldb"
	"github.com/onorbit/pixelite/image"
	"github.com/onorbit/pixelite/pkg/log"
)

var ErrThumbnailNotAvailable = errors.New("thumbnail is not available")

type thumbnailedAlbum struct {
	thumbnailedAlbumID  int64
	albumIDHash         string
	createTimestamp     time.Time
	lastAccessTimestamp time.Time
	thumbnailDir        string
	thumbnails          map[string]string
	progress            map[string]*sync.Cond
	mutex               sync.Mutex
}

func newThumbnailedAlbum(databaseID int64, albumID string, createTimestamp, lastAccessTimestamp time.Time, thumbnailLibDir string, makeDir bool) *thumbnailedAlbum {
	albumIDHash := getAlbumIDHash(albumID)
	leafThumbnailDir := fmt.Sprintf("%s_%x", albumIDHash, createTimestamp.Unix())
	thumbnailDir := filepath.Join(thumbnailLibDir, albumIDHash[0:2], albumIDHash[2:4], leafThumbnailDir)

	if makeDir {
		if err := os.MkdirAll(thumbnailDir, 0700); err != nil {
			log.Error("failed to make thumbnail path [%s] - [%v]", thumbnailDir, err.Error())
			return nil
		}
	}

	albumInfo := &thumbnailedAlbum{
		thumbnailedAlbumID:  databaseID,
		albumIDHash:         albumIDHash,
		createTimestamp:     createTimestamp,
		lastAccessTimestamp: lastAccessTimestamp,
		thumbnailDir:        thumbnailDir,
		thumbnails:          make(map[string]string),
		progress:            make(map[string]*sync.Cond),
	}

	return albumInfo
}

func (a *thumbnailedAlbum) getThumbnailPath(albumPath, fileName string) (string, error) {
	origImgPath := path.Join(albumPath, fileName)
	thumbnailName := strings.TrimSuffix(fileName, filepath.Ext(fileName)) + ".jpg"

	a.mutex.Lock()
	defer a.mutex.Unlock()

	// thumbnail already exists. return it directly.
	thumbnailPath, ok := a.thumbnails[origImgPath]
	if ok {
		return thumbnailPath, nil
	}

	// check if the image is already being processed.
	var cond *sync.Cond
	if cond, ok = a.progress[origImgPath]; !ok {
		cond = sync.NewCond(&a.mutex)
		a.progress[origImgPath] = cond

		destPath := filepath.Join(a.thumbnailDir, thumbnailName)
		go a.buildThumbnail(origImgPath, destPath, cond)
	}

	// wait for the thumbnail to be ready.
	cond.Wait()

	// thumbnail path could be not registered, if there was some error on building.
	thumbnailPath, ok = a.thumbnails[origImgPath]
	if !ok {
		return "", ErrThumbnailNotAvailable
	}

	return thumbnailPath, nil
}

func (a *thumbnailedAlbum) buildThumbnail(origImgPath, thumbnailPath string, cond *sync.Cond) {
	// get parameters for making thumbnail.
	thumbnailDim := config.Get().Thumbnail.MaxDimension
	thumbnailJpegQuality := config.Get().Thumbnail.JpegQuality

	// make actual thumbnail.
	err := image.MakeThumbnail(origImgPath, thumbnailPath, thumbnailDim, thumbnailJpegQuality, true)
	if err != nil {
		log.Error("failed to make thumbnail image for [%s] - [%v]", origImgPath, err.Error())
		cond.Broadcast()
		return
	}

	a.mutex.Lock()
	delete(a.progress, origImgPath)
	a.thumbnails[origImgPath] = thumbnailPath
	a.mutex.Unlock()

	// TODO : check if there could be some races between deleting flow and this.
	globaldb.InsertThumbnail(origImgPath, thumbnailPath, a.thumbnailedAlbumID)

	cond.Broadcast()
}

func (a *thumbnailedAlbum) cleanUp() {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// delete thumbnail files.
	for _, thumbnailPath := range a.thumbnails {
		if err := os.Remove(thumbnailPath); err != nil {
			log.Error("failed to remove thumbnail file [%s] - [%v]", thumbnailPath, err.Error())
		}
	}

	// delete thumbnail directory.
	if err := os.Remove(a.thumbnailDir); err != nil {
		log.Error("failed to remove thumbnail directory [%s] - [%v]", a.thumbnailDir, err.Error())
	}
}
