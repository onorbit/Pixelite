package thumbnail

import (
	"os"
	"path/filepath"

	"github.com/onorbit/pixelite/config"
	"github.com/onorbit/pixelite/media"
	"github.com/onorbit/pixelite/pkg/log"
)

func makeCover(fileName, albumPath, albumID, libraryID string) (string, error) {
	coverCfg := config.Get().Cover

	// TODO : following is called repeatedly.
	libraryPath := filepath.Join(coverCfg.StorePath, libraryID)
	if err := os.MkdirAll(libraryPath, 0700); err != nil {
		log.Error("failed to make cover path [%s] - [%v]", libraryPath, err.Error())
		return "", err
	}

	coverFileName := getAlbumIDHash(albumID) + ".jpg"
	coverPath := filepath.Join(coverCfg.StorePath, libraryID, coverFileName)
	origImgPath := filepath.Join(albumPath, fileName)

	// make actual cover.
	mediaFile, err := media.LoadMediaFile(origImgPath)
	if err != nil {
		log.Error("failed to load media file for cover [%s] - [%v]", origImgPath, err.Error())
		return "", err
	}

	err = mediaFile.MakeThumbnail(coverPath, coverCfg.MaxDimension, coverCfg.JpegQuality, true)
	if err != nil {
		log.Error("failed to make cover image for [%s] - [%v]", origImgPath, err.Error())
		return "", err
	}

	return coverPath, nil
}
