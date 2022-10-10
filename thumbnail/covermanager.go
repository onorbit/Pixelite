package thumbnail

import (
	"sync"

	"github.com/onorbit/pixelite/database/globaldb"
)

type coverManager struct {
	covers map[thumbnailedAlbumKey]string
	mutex  sync.Mutex
}

var gCoverManager coverManager

func (m *coverManager) getAlbumCover(fileName, albumPath, albumID, libraryID string) string {
	albumKey := thumbnailedAlbumKey{
		libraryID: libraryID,
		albumID:   albumID,
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	coverPath, ok := m.covers[albumKey]
	if ok {
		// cover exists. return directly.
		return coverPath
	}

	// TODO : use progress conditionals to prevent duplicated works.
	coverPath, err := makeCover(fileName, albumPath, albumID, libraryID)
	if err != nil {
		return ""
	}

	globaldb.InsertAlbumCover(libraryID, albumID, coverPath)
	m.covers[albumKey] = coverPath

	return coverPath
}

func (m *coverManager) loadFromGlobalDB() error {
	coverRows, err := globaldb.LoadAllAlbumCovers()
	if err != nil {
		return err
	}

	for _, row := range coverRows {
		albumKey := thumbnailedAlbumKey{
			libraryID: row.LibraryID,
			albumID:   row.AlbumID,
		}

		m.covers[albumKey] = row.CoverPath
	}

	return nil
}
