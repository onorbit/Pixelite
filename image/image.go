package image

import (
	"errors"
	"image"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/chai2010/webp"
	"github.com/rwcarlsen/goexif/exif"
)

var imageExt = []string{
	".jpg",
	".png",
	".webp",
}

var ErrFormatNotSupported = errors.New("not supported file format")

func openImage(srcPath string) (img image.Image, imgExif *exif.Exif, err error) {
	imgFile, err := os.Open(srcPath)
	if err != nil {
		return
	}
	defer imgFile.Close()

	// attempt to extract EXIF information.
	imgExif, _ = exif.Decode(imgFile)
	imgFile.Seek(0, io.SeekStart)

	// decode image according to the extension.
	fileExt := strings.ToLower(filepath.Ext(srcPath))
	switch fileExt {
	case ".jpg":
		fallthrough
	case ".png":
		img, _, err = image.Decode(imgFile)
		return
	case ".webp":
		img, err = webp.Decode(imgFile)
		return
	default:
		err = ErrFormatNotSupported
		return
	}
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

func MakeThumbnail(srcPath, dstPath string, thumbnailSize, jpegQuality int, squareCrop bool) error {
	workImage, imgExif, err := openImage(srcPath)
	if err != nil {
		return err
	}

	// rotate if necessary.
	if imgExif != nil {
		orientationTag, err := imgExif.Get(exif.Orientation)
		if !exif.IsTagNotPresentError(err) && orientationTag != nil {
			orientation, _ := orientationTag.Int(0)
			switch orientation {
			case 2:
				workImage = transform.FlipH(workImage)
			case 3:
				workImage = transform.Rotate(workImage, 180.0, &transform.RotationOptions{ResizeBounds: true})
			case 4:
				workImage = transform.FlipV(workImage)
			case 6:
				workImage = transform.Rotate(workImage, 90, &transform.RotationOptions{ResizeBounds: true})
			case 8:
				workImage = transform.Rotate(workImage, 270, &transform.RotationOptions{ResizeBounds: true})
			}
		}
	}

	bounds := workImage.Bounds()
	imageWidth := bounds.Dx()
	imageHeight := bounds.Dy()

	// crop to square, if necessary.
	if squareCrop && imageWidth != imageHeight {
		cropRect := image.Rectangle{}
		if imageWidth > imageHeight {
			cropRect.Min.Y = 0
			cropRect.Max.Y = imageHeight

			cropRect.Min.X = (imageWidth - imageHeight) / 2
			cropRect.Max.X = cropRect.Min.X + imageHeight

			imageWidth = imageHeight
		} else {
			cropRect.Min.X = 0
			cropRect.Max.X = imageWidth

			cropRect.Min.Y = (imageHeight - imageWidth) / 2
			cropRect.Max.Y = cropRect.Min.Y + imageWidth

			imageHeight = imageWidth
		}

		workImage = transform.Crop(workImage, cropRect)
	}

	// determine thumbnail size.
	thumbnailWidth := 0
	thumbnailHeight := 0
	if imageWidth > imageHeight && imageWidth > thumbnailSize {
		thumbnailWidth = thumbnailSize
		thumbnailHeight = int(float32(imageHeight) * (float32(thumbnailSize) / float32(imageWidth)))
	} else if imageHeight > thumbnailSize {
		thumbnailWidth = int(float32(imageWidth) * (float32(thumbnailSize) / float32(imageHeight)))
		thumbnailHeight = thumbnailSize
	}

	if thumbnailWidth == 0 && thumbnailHeight == 0 {
		err = imgio.Save(dstPath, workImage, imgio.JPEGEncoder(jpegQuality))
	} else {
		thumbnail := transform.Resize(workImage, thumbnailWidth, thumbnailHeight, transform.CatmullRom)
		err = imgio.Save(dstPath, thumbnail, imgio.JPEGEncoder(jpegQuality))
	}

	return err
}
