package media

import (
	"errors"
	"path/filepath"
	"strings"
)

var ErrFormatNotSupported = errors.New("not supported file format")

type MediaFile interface {
	MakeThumbnail(dstPath string, thumbnailSize, jpegQuality int, squareCrop bool) error
}

type MediaFileLoader func(srcPath string) (MediaFile, error)

var gMediaFileLoaders map[string]MediaFileLoader

func Initialize() {
	gMediaFileLoaders = make(map[string]MediaFileLoader)

	registerImageLoaders()
}

func IsSupportedMedia(fileName string) bool {
	fileExt := filepath.Ext(fileName)
	fileExt = strings.ToLower(fileExt)

	_, ok := gMediaFileLoaders[fileExt]
	return ok
}

func LoadMediaFile(srcPath string) (MediaFile, error) {
	fileExt := filepath.Ext(srcPath)
	fileExt = strings.ToLower(fileExt)

	loader, ok := gMediaFileLoaders[fileExt]
	if !ok {
		return nil, ErrFormatNotSupported
	}

	return loader(srcPath)
}
