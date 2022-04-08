package thumbnail

import (
	"errors"
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
	lastAccessTimestamp time.Time
	thumbnails          map[string]string
	progress            map[string]*sync.Cond
	mutex               sync.Mutex
}

func newThumbnailedAlbum(databaseID int64, albumID string, lastAccessTimestamp time.Time) *thumbnailedAlbum {
	albumInfo := &thumbnailedAlbum{
		thumbnailedAlbumID:  databaseID,
		albumIDHash:         getAlbumIDHash(albumID),
		lastAccessTimestamp: lastAccessTimestamp,
		thumbnails:          make(map[string]string),
		progress:            make(map[string]*sync.Cond),
	}

	return albumInfo
}

func (a *thumbnailedAlbum) getThumbnailPath(thumbnailLibPath, albumPath, fileName string) (string, error) {
	origImgPath := path.Join(albumPath, fileName)
	thumbnailName := strings.TrimSuffix(fileName, filepath.Ext(fileName)) + ".jpg"

	a.mutex.Lock()
	defer a.mutex.Unlock()

	// thumbnail already exists. return it directly.
	thumbnailPath, ok := a.thumbnails[origImgPath]
	if ok {
		return thumbnailPath, nil
	}

	// prepare thumbnail path for the album.
	albumIDHash := a.albumIDHash
	thumbnailDir := filepath.Join(thumbnailLibPath, albumIDHash[0:2], albumIDHash[2:4], albumIDHash)

	// check if the image is already being processed.
	var cond *sync.Cond
	if cond, ok = a.progress[origImgPath]; !ok {
		cond = sync.NewCond(&a.mutex)
		a.progress[origImgPath] = cond

		destPath := filepath.Join(thumbnailDir, thumbnailName)
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
	err := image.MakeThumbnail(origImgPath, thumbnailPath, thumbnailDim, thumbnailJpegQuality)
	if err != nil {
		log.Error("failed to make thumbnail image for [%s] - [%v]", origImgPath, err.Error())
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
