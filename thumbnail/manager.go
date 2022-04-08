package thumbnail

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/onorbit/pixelite/config"
	"github.com/onorbit/pixelite/database/globaldb"
	"github.com/onorbit/pixelite/pkg/log"
)

const tickIntervalMinutes = 5

type thumbnailedAlbumKey struct {
	libraryID string
	albumID   string
}

type manager struct {
	progress           map[string]*sync.Cond
	thumbnailedAlbums  map[thumbnailedAlbumKey]*thumbnailedAlbum
	recentAccessAlbums map[thumbnailedAlbumKey]struct{}
	cancelTickerFunc   context.CancelFunc
	tickerWaitGroup    sync.WaitGroup
	mutex              sync.Mutex
}

var gManager manager

func getAlbumIDHash(albumID string) string {
	albumIDHashArr := md5.Sum([]byte(albumID))
	return hex.EncodeToString(albumIDHashArr[:])
}

func (m *manager) getThumbnailPath(fileName, albumPath, albumID, libraryID string) string {
	thumbnailLibDir := filepath.Join(config.Get().Thumbnail.StorePath, libraryID)
	m.mutex.Lock()

	// update access timestamp per album.
	currTime := time.Now()
	albumKey := thumbnailedAlbumKey{
		libraryID: libraryID,
		albumID:   albumID,
	}

	albumInfo, ok := m.thumbnailedAlbums[albumKey]
	if ok {
		albumInfo.lastAccessTimestamp = currTime
		m.recentAccessAlbums[albumKey] = struct{}{}
	} else {
		// insert to global db.
		databaseID, err := globaldb.InsertThumbnailedAlbum(libraryID, albumID, currTime, currTime)
		if err != nil {
			log.Error("failed to insert thumbnailed album - libraryID [%s], albumID [%s] - %v", libraryID, albumID, err.Error())
			m.mutex.Unlock()
			return ""
		}

		albumInfo = newThumbnailedAlbum(databaseID, albumID, currTime)

		// prepare thumbnail file path.
		albumIDHash := albumInfo.albumIDHash
		thumbnailDir := filepath.Join(thumbnailLibDir, albumIDHash[0:2], albumIDHash[2:4], albumIDHash)
		if err := os.MkdirAll(thumbnailDir, 0700); err != nil {
			log.Error("failed to make thumbnail path [%s] - [%v]", thumbnailDir, err.Error())
			m.mutex.Unlock()
			return ""
		}

		// register the structure.
		m.thumbnailedAlbums[albumKey] = albumInfo
		log.Info("album [%s] - [%s] in library [%s] is registered as thumbnailed", albumID, albumInfo.albumIDHash, libraryID)
	}

	m.mutex.Unlock()

	ret, err := albumInfo.getThumbnailPath(thumbnailLibDir, albumPath, fileName)
	if err != nil {
		// TODO : add parameters to following log
		log.Error("failed to get thumbnail path - [%v]", err)
		return ""
	}

	return ret
}

func (m *manager) startTick() {
	tickFunc := func(ctx context.Context, wg *sync.WaitGroup) {
		ticker := time.NewTicker(tickIntervalMinutes * time.Minute)

		for {
			select {
			case <-ctx.Done():
				wg.Done()
				gManager.syncLastAccessTimeToDB()
				return
			case <-ticker.C:
				gManager.syncLastAccessTimeToDB()
				// TODO : delete thumbnails.
			}
		}
	}

	m.tickerWaitGroup.Add(1)

	ctx, cancel := context.WithCancel(context.Background())
	go tickFunc(ctx, &m.tickerWaitGroup)
	m.cancelTickerFunc = cancel

	log.Info("thumbnail manager starts ticking")
}

func (m *manager) stopTick() {
	m.cancelTickerFunc()
	m.tickerWaitGroup.Wait()

	log.Info("thumbnail manager stopped ticking")
}

func (m *manager) syncLastAccessTimeToDB() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for key, _ := range m.recentAccessAlbums {
		albumInfo := m.thumbnailedAlbums[key]
		globaldb.UpdateThumbnailedAlbumAccessTimestamp(albumInfo.thumbnailedAlbumID, albumInfo.lastAccessTimestamp)

		delete(m.recentAccessAlbums, key)
	}
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
		progress:           make(map[string]*sync.Cond),
		thumbnailedAlbums:  make(map[thumbnailedAlbumKey]*thumbnailedAlbum),
		recentAccessAlbums: make(map[thumbnailedAlbumKey]struct{}),
		mutex:              sync.Mutex{},
	}

	// load thumbnailed albums from global DB.
	thumbnailedAlbumRows, err := globaldb.LoadAllThumbnailedAlbums()
	if err != nil {
		return err
	}

	// temporary index for initial loading.
	albumsByDBID := make(map[int64]*thumbnailedAlbum)

	for _, row := range thumbnailedAlbumRows {
		entry := newThumbnailedAlbum(row.ID, getAlbumIDHash(row.AlbumID), time.Unix(row.LastAccessTimestamp, 0))
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

	gManager.startTick()

	return nil
}

func Cleanup() {
	gManager.stopTick()
}

func GetThumbnailPath(fileName, albumPath, albumID, libraryID string) string {
	return gManager.getThumbnailPath(fileName, albumPath, albumID, libraryID)
}
