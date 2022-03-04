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
	thumbnailedAlbumID  int64
	albumIDHash         string
	lastAccessTimestamp time.Time
	thumbnails          map[string]string
}

type thumbnailedAlbumKey struct {
	libraryID string
	albumID   string
}

type manager struct {
	//thumbnails        map[string]string // this may be moved into thumbnailedAlbum
	progress          map[string]*sync.Cond
	thumbnailedAlbums map[thumbnailedAlbumKey]*thumbnailedAlbum
	mutex             sync.Mutex
}

var gManager manager

func getAlbumIDHash(albumID string) string {
	albumIDHashArr := md5.Sum([]byte(albumID))
	return hex.EncodeToString(albumIDHashArr[:])
}

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
		albumInfo.lastAccessTimestamp = currTime
	} else {
		databaseID, err := globaldb.InsertThumbnailedAlbum(libraryID, albumID, currTime)
		if err != nil {
			// TODO : handle the error properly.
		}

		albumInfo = &thumbnailedAlbum{
			thumbnailedAlbumID:  databaseID,
			albumIDHash:         getAlbumIDHash(albumID),
			lastAccessTimestamp: currTime,
			thumbnails:          make(map[string]string),
		}
		m.thumbnailedAlbums[albumKey] = albumInfo
	}

	// prepare thumbnail file path.
	albumIDHash := albumInfo.albumIDHash
	thumbnailDir := filepath.Join(thumbnailLibDir, albumIDHash[0:2], albumIDHash[2:4], albumIDHash)
	if !ok {
		if err := os.MkdirAll(thumbnailDir, 0700); err != nil {
			// TODO : handle the error properly.
			return ""
		}
	}

	// thumbnail already exists. return it directly.
	existThumbnailPath, ok := albumInfo.thumbnails[origImgPath]
	if ok {
		return existThumbnailPath
	}

	// check if the image is already being processed.
	var cond *sync.Cond
	if cond, ok = m.progress[origImgPath]; !ok {
		cond = sync.NewCond(&m.mutex)
		m.progress[origImgPath] = cond

		thumbnailPath := filepath.Join(thumbnailDir, thumbnailName)
		go m.buildThumbnail(albumKey, albumInfo.thumbnailedAlbumID, origImgPath, thumbnailPath, cond)
	}

	// wait for the thumbnail to be made.
	cond.Wait()

	// at this point, there's no guarantee that album info and thumbnail info exist.
	albumInfo, ok = m.thumbnailedAlbums[albumKey]
	if !ok {
		return ""
	}

	thumbnailPath, ok := albumInfo.thumbnails[origImgPath]
	if !ok {
		return ""
	}

	return thumbnailPath
}

func (m *manager) buildThumbnail(thumbnailedAlbumKey thumbnailedAlbumKey, thumbnailedAlbumID int64, imgPath, thumbnailPath string, signalCond *sync.Cond) {
	// get parameters for making thumbnail.
	thumbnailDim := config.Get().Thumbnail.MaxDimension
	thumbnailJpegQuality := config.Get().Thumbnail.JpegQuality

	// make actual thumbnail.
	err := image.MakeThumbnail(imgPath, thumbnailPath, thumbnailDim, thumbnailJpegQuality)

	if err == nil {
		globaldb.RegisterThumbnail(imgPath, thumbnailPath, thumbnailedAlbumID)
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.progress, imgPath)
	if err == nil {
		if albumInfo, ok := m.thumbnailedAlbums[thumbnailedAlbumKey]; ok {
			albumInfo.thumbnails[imgPath] = thumbnailPath
		} else {
			// TODO : delete thumbnail file.
		}
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
		progress:          make(map[string]*sync.Cond),
		thumbnailedAlbums: make(map[thumbnailedAlbumKey]*thumbnailedAlbum),
		mutex:             sync.Mutex{},
	}

	// load thumbnailed albums from global DB.
	thumbnailedAlbumRows, err := globaldb.LoadAllThumbnailedAlbums()
	if err != nil {
		return err
	}

	// temporary index for initial loading.
	albumsByDBID := make(map[int64]*thumbnailedAlbum)

	for _, row := range thumbnailedAlbumRows {
		entry := &thumbnailedAlbum{
			thumbnailedAlbumID:  row.ID,
			albumIDHash:         getAlbumIDHash(row.AlbumID),
			lastAccessTimestamp: time.Unix(row.LastAccessTimestamp, 0),
			thumbnails:          make(map[string]string),
		}

		key := thumbnailedAlbumKey{
			libraryID: row.LibraryID,
			albumID:   row.AlbumID,
		}

		gManager.thumbnailedAlbums[key] = entry

		// index the entry for further loading.
		albumsByDBID[row.ID] = entry
	}

	// load thumbnails from global DB.
	thumbnailRows, err := globaldb.LoadAllThumbnails()
	if err != nil {
		return err
	}

	for _, row := range thumbnailRows {
		thumbnailedAlbum, ok := albumsByDBID[row.ThumbnailedAlbumID]
		if !ok {
			// TODO : remove the row and thumbnail file.
			continue
		}

		thumbnailedAlbum.thumbnails[row.ImagePath] = row.ThumbnailPath
	}

	return nil
}

func GetThumbnailPath(fileName, albumPath, albumID, libraryID string) string {
	return gManager.getThumbnailPath(fileName, albumPath, albumID, libraryID)
}
