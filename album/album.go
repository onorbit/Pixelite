package album

import (
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/onorbit/pixelite/media"
	"github.com/onorbit/pixelite/pkg/fileutils"
)

type Album struct {
	id                   string
	path                 string
	coverFileName        string
	mediaListCache       []string
	mediaListCacheExpire time.Time
}

const mediaListCacheLifetime = time.Hour * 24

func NewAlbum(id, path, coverFileName string) Album {
	newAlbum := Album{
		id:            id,
		path:          path,
		coverFileName: coverFileName,
	}

	return newAlbum
}

func (a Album) GetID() string {
	return a.id
}

func (a Album) GetPath() string {
	return a.path
}

func (a Album) GetCoverFileName() string {
	return a.coverFileName
}

func (a Album) ListMedias() ([]string, error) {
	if a.mediaListCache != nil && time.Now().Before(a.mediaListCacheExpire) {
		return a.mediaListCache, nil
	}

	content, err := ioutil.ReadDir(a.path)
	if err != nil {
		return nil, err
	}

	mediaList := make([]string, 0, len(content))
	for _, entry := range content {
		if entry.IsDir() {
			continue
		}

		if !media.IsSupportedMedia(entry.Name()) {
			continue
		}

		filePath := filepath.Join(a.path, entry.Name())
		if isHidden, err := fileutils.IsHidden(filePath); err != nil {
			// TODO : log?
			continue
		} else if isHidden {
			continue
		}

		mediaList = append(mediaList, entry.Name())
	}

	a.mediaListCache = mediaList
	a.mediaListCacheExpire = time.Now().Add(mediaListCacheLifetime)

	return mediaList, nil
}
