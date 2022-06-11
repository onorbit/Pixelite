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
	thumbnailedAlbums  map[thumbnailedAlbumKey]*thumbnailedAlbum
	recentAccessAlbums map[thumbnailedAlbumKey]struct{}
	albumCovers        map[thumbnailedAlbumKey]string
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

	currTime := time.Now()
	albumKey := thumbnailedAlbumKey{
		libraryID: libraryID,
		albumID:   albumID,
	}

	m.mutex.Lock()

	albumInfo, ok := m.thumbnailedAlbums[albumKey]
	if ok {
		// update access timestamp per album.
		albumInfo.lastAccessTimestamp = currTime
		m.recentAccessAlbums[albumKey] = struct{}{}
	} else {
		// insert to global db.
		databaseID, err := globaldb.InsertThumbnailedAlbum(libraryID, albumID)
		if err != nil {
			log.Error("failed to insert thumbnailed album - libraryID [%s], albumID [%s] - %v", libraryID, albumID, err.Error())
			m.mutex.Unlock()
			return ""
		}

		albumInfo = newThumbnailedAlbum(databaseID, albumID, currTime, currTime, thumbnailLibDir, true)
		if albumInfo == nil {
			// failed to create thumbnail directory.
			m.mutex.Unlock()
			return ""
		}

		// register the structure.
		m.thumbnailedAlbums[albumKey] = albumInfo
		log.Info("album [%s] - [%s] in library [%s] is registered to thumbnail manager", albumID, albumInfo.albumIDHash, libraryID)
	}

	m.mutex.Unlock()

	ret, err := albumInfo.getThumbnailPath(albumPath, fileName)
	if err != nil {
		// TODO : add parameters to following log
		log.Error("failed to get thumbnail path - [%v]", err)
		return ""
	}

	return ret
}

func (m *manager) getAlbumCover(fileName, albumPath, albumID, libraryID string) string {
	albumKey := thumbnailedAlbumKey{
		libraryID: libraryID,
		albumID:   albumID,
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	coverPath, ok := m.albumCovers[albumKey]
	if ok {
		// cover exists. return directly.
		return coverPath
	}

	// TODO : use progress conditionals to prevent duplicated works.
	coverPath, err := makeCover(fileName, albumPath, albumID, libraryID)
	if err != nil {
		return ""
	}

	m.albumCovers[albumKey] = coverPath
	return coverPath
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
				gManager.deleteUnusedThumbnails()
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

func (m *manager) deleteUnusedThumbnails() {
	thresholdTime := time.Now().Add(time.Hour * 24 * time.Duration(config.Get().Thumbnail.LifetimeUnusedDays) * -1)
	toDelete := make([]*thumbnailedAlbum, 0)

	// select thumbnailedAlbum to delete with lock acquisition.
	m.mutex.Lock()
	for key, albumInfo := range m.thumbnailedAlbums {
		if albumInfo.lastAccessTimestamp.Before(thresholdTime) {
			toDelete = append(toDelete, albumInfo)

			// TODO : try deleting the DB entries outside of lock scope.
			globaldb.DeleteThumbnailedAlbum(albumInfo.thumbnailedAlbumID)
			delete(m.thumbnailedAlbums, key)
			delete(m.recentAccessAlbums, key)

			log.Info("album [%s] in library [%s] is unregistered from thumbnail manager", key.albumID, key.libraryID)
		}
	}
	m.mutex.Unlock()

	// delete thumbnail files and directory.
	for _, albumInfo := range toDelete {
		albumInfo.cleanUp()
	}
}

//-----------------------------------------------------------------------------
// public functions
//-----------------------------------------------------------------------------

func Initialize() error {
	cfg := config.Get()

	// make paths.
	if err := os.MkdirAll(cfg.Thumbnail.StorePath, 0700); err != nil {
		return err
	}
	if err := os.MkdirAll(cfg.Cover.StorePath, 0700); err != nil {
		return err
	}

	// initialize manager.
	gManager = manager{
		thumbnailedAlbums:  make(map[thumbnailedAlbumKey]*thumbnailedAlbum),
		recentAccessAlbums: make(map[thumbnailedAlbumKey]struct{}),
		albumCovers:        make(map[thumbnailedAlbumKey]string),
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
		thumbnailLibDir := filepath.Join(cfg.Thumbnail.StorePath, row.LibraryID)
		createTimestamp := time.Unix(row.CreateTimestamp, 0)
		lastAccessTimestamp := time.Unix(row.LastAccessTimestamp, 0)

		entry := newThumbnailedAlbum(row.ID, row.AlbumID, createTimestamp, lastAccessTimestamp, thumbnailLibDir, false)
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

func GetAlbumCover(fileName, albumPath, albumID, libraryID string) string {
	return gManager.getAlbumCover(fileName, albumPath, albumID, libraryID)
}
