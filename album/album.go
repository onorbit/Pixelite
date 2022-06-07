package album

import (
	"io/ioutil"
	"path/filepath"

	"github.com/onorbit/pixelite/image"
	"github.com/onorbit/pixelite/pkg/fileutils"
)

type Album struct {
	id            string
	path          string
	coverFileName string
}

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

func (a Album) ListImages() ([]string, error) {
	content, err := ioutil.ReadDir(a.path)
	if err != nil {
		return nil, err
	}

	imageList := make([]string, 0, len(content))
	for _, entry := range content {
		if entry.IsDir() {
			continue
		}

		if !image.IsImageFile(entry.Name()) {
			continue
		}

		filePath := filepath.Join(a.path, entry.Name())
		if isHidden, err := fileutils.IsHidden(filePath); err != nil {
			// TODO : log?
			continue
		} else if isHidden {
			continue
		}

		imageList = append(imageList, entry.Name())
	}

	return imageList, nil
}
