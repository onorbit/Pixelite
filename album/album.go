package album

import (
	"io/ioutil"

	"github.com/onorbit/pixelite/image"
)

type Album struct {
	path string
}

func NewAlbum(path string) Album {
	newAlbum := Album{
		path: path,
	}

	return newAlbum
}

func (a Album) GetPath() string {
	return a.path
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

		if image.IsImageFile(entry.Name()) {
			imageList = append(imageList, entry.Name())
		}
	}

	return imageList, nil
}
