package thumbnail

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/onorbit/pixelite/config"
	"github.com/onorbit/pixelite/database/globaldb"
	"github.com/onorbit/pixelite/image"
)

type thumbnailedAlbum struct {
	isDirty             bool
	albumIDHash         string
	lastAccessTimestamp time.Time
}

type thumbnailedAlbumKey struct {
	libraryID string
	albumID   string
}

type manager struct {
	thumbnails        map[string]string // this may be moved into thumbnailedAlbum
	progress          map[string]*sync.Cond
	thumbnailedAlbums map[thumbnailedAlbumKey]*thumbnailedAlbum
	mutex             sync.Mutex
}

var gManager manager

func (m *manager) getThumbnailPath(fileName, albumPath, albumID, libraryID string) string {
	// prepare path elements outside of mutex scope.
	thumbnailLibDir := filepath.Join(config.Get().Thumbnail.StorePath, libraryID)
	thumbnailName := strings.TrimSuffix(fileName, filepath.Ext(fileName)) + ".jpg"
	origImgPath := path.Join(albumPath, fileName)

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// update access timestamp per album.
	currTime := time.Now()
	albumKey := thumbnailedAlbumKey{
		libraryID: libraryID,
		albumID:   albumID,
	}

	albumInfo, ok := m.thumbnailedAlbums[albumKey]
	if ok {
		albumInfo.isDirty = true
		albumInfo.lastAccessTimestamp = currTime
	} else {
		albumIDHashArr := md5.Sum([]byte(albumID))
		albumIDHash := hex.EncodeToString(albumIDHashArr[:])

		albumInfo = &thumbnailedAlbum{
			isDirty:             true,
			albumIDHash:         albumIDHash,
			lastAccessTimestamp: currTime,
		}
		m.thumbnailedAlbums[albumKey] = albumInfo
	}

	// prepare thumbnail file path.
	thumbnailDir := filepath.Join(thumbnailLibDir, albumInfo.albumIDHash)
	if !ok {
		if err := os.MkdirAll(thumbnailDir, 0700); err != nil {
			// TODO : handle the error properly.
			return ""
		}
	}

	thumbnailPath := filepath.Join(thumbnailDir, thumbnailName)

	// thumbnail already exists. return it directly.
	existThumbnailPath, ok := m.thumbnails[origImgPath]
	if ok {
		return existThumbnailPath
	}

	var cond *sync.Cond
	if cond, ok = m.progress[origImgPath]; !ok {
		cond = sync.NewCond(&m.mutex)
		m.progress[origImgPath] = cond

		go m.buildThumbnail(origImgPath, thumbnailPath, cond)
	}

	cond.Wait()

	thumbnailPath, ok = m.thumbnails[origImgPath]
	if !ok {
		return ""
	}

	return thumbnailPath
}

func (m *manager) buildThumbnail(imgPath, thumbnailPath string, signalCond *sync.Cond) {
	// get parameters for making thumbnail.
	thumbnailDim := config.Get().Thumbnail.MaxDimension
	thumbnailJpegQuality := config.Get().Thumbnail.JpegQuality

	// make actual thumbnail.
	err := image.MakeThumbnail(imgPath, thumbnailPath, thumbnailDim, thumbnailJpegQuality)

	if err == nil {
		globaldb.RegisterThumbnail(imgPath, thumbnailPath)
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.progress, imgPath)
	if err == nil {
		m.thumbnails[imgPath] = thumbnailPath
	}
	signalCond.Broadcast()
}

//-----------------------------------------------------------------------------
// public functions
//-----------------------------------------------------------------------------

func Initialize() error {
	thumbnailStorePath := config.Get().Thumbnail.StorePath
	if err := os.MkdirAll(thumbnailStorePath, 0700); err != nil {
		return err
	}

	gManager = manager{
		thumbnails:        make(map[string]string),
		progress:          make(map[string]*sync.Cond),
		thumbnailedAlbums: make(map[thumbnailedAlbumKey]*thumbnailedAlbum),
		mutex:             sync.Mutex{},
	}

	thumbnailRows, err := globaldb.LoadAllThumbnails()
	if err != nil {
		return err
	}

	for _, row := range thumbnailRows {
		gManager.thumbnails[row.ImagePath] = row.ThumbnailPath
	}

	return nil
}

func GetThumbnailPath(fileName, albumPath, albumID, libraryID string) string {
	return gManager.getThumbnailPath(fileName, albumPath, albumID, libraryID)
}
