package image

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

var imageExt = []string{
	".jpg",
	".png",
}

func IsImageFile(fileName string) bool {
	fileExt := filepath.Ext(fileName)
	fileExt = strings.ToLower(fileExt)
	for _, ext := range imageExt {
		if ext == fileExt {
			return true
		}
	}

	return false
}

func MakeThumbnail(srcPath, dstPath string, thumbnailSize, jpegQuality int) error {
	image, err := imgio.Open(srcPath)
	if err != nil {
		return err
	}

	// determine thumbnail size.
	bounds := image.Bounds()
	imageWidth := bounds.Dx()
	imageHeight := bounds.Dy()

	thumbnailWidth := 0
	thumbnailHeight := 0
	if imageWidth > imageHeight && imageWidth > thumbnailSize {
		thumbnailWidth = thumbnailSize
		thumbnailHeight = int(float32(imageHeight) * (float32(thumbnailSize) / float32(imageWidth)))
	} else if imageHeight > thumbnailSize {
		thumbnailWidth = int(float32(imageWidth) * (float32(thumbnailSize) / float32(imageHeight)))
		thumbnailHeight = thumbnailSize
	}

	// serve the original image directly, as it is small enough.
	if thumbnailWidth == 0 && thumbnailHeight == 0 {
		srcFile, err := os.Open(srcPath)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.Open(dstPath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		return err
	}

	// perform resize.
	thumbnail := transform.Resize(image, thumbnailWidth, thumbnailHeight, transform.CatmullRom)
	err = imgio.Save(dstPath, thumbnail, imgio.JPEGEncoder(jpegQuality))

	return err
}
