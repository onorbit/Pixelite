package media

import (
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

func imageFileLoader(srcPath string) (MediaFile, error) {
	imgFile, err := os.Open(srcPath)
	if err != nil {
		return nil, err
	}
	defer imgFile.Close()

	// attempt to extract EXIF information.
	imgExif, _ := exif.Decode(imgFile)
	imgFile.Seek(0, io.SeekStart)

	// decode image according to the extension.
	var imageData image.Image
	fileExt := strings.ToLower(filepath.Ext(srcPath))
	switch fileExt {
	case ".jpg":
		fallthrough
	case ".jpeg":
		fallthrough
	case ".png":
		imageData, _, err = image.Decode(imgFile)
	case ".webp":
		imageData, err = webp.Decode(imgFile)
	default:
		err = ErrFormatNotSupported
	}

	if err != nil {
		return nil, err
	}

	ret := &imageFile{
		imageData: imageData,
		exifData:  imgExif,
	}

	return ret, nil
}

func registerImageLoaders() {
	gMediaFileLoaders[".jpg"] = imageFileLoader
	gMediaFileLoaders[".jpeg"] = imageFileLoader
	gMediaFileLoaders[".png"] = imageFileLoader
	gMediaFileLoaders[".webp"] = imageFileLoader
}

type imageFile struct {
	imageData image.Image
	exifData  *exif.Exif
}

func (i *imageFile) MakeThumbnail(dstPath string, thumbnailSize, jpegQuality int, squareCrop bool) error {
	// rotate if necessary.
	if i.exifData != nil {
		orientationTag, err := i.exifData.Get(exif.Orientation)
		if !exif.IsTagNotPresentError(err) && orientationTag != nil {
			orientation, _ := orientationTag.Int(0)
			switch orientation {
			case 2:
				i.imageData = transform.FlipH(i.imageData)
			case 3:
				i.imageData = transform.Rotate(i.imageData, 180.0, &transform.RotationOptions{ResizeBounds: true})
			case 4:
				i.imageData = transform.FlipV(i.imageData)
			case 6:
				i.imageData = transform.Rotate(i.imageData, 90, &transform.RotationOptions{ResizeBounds: true})
			case 8:
				i.imageData = transform.Rotate(i.imageData, 270, &transform.RotationOptions{ResizeBounds: true})
			}
		}
	}

	bounds := i.imageData.Bounds()
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

		i.imageData = transform.Crop(i.imageData, cropRect)
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

	var err error
	if thumbnailWidth == 0 && thumbnailHeight == 0 {
		err = imgio.Save(dstPath, i.imageData, imgio.JPEGEncoder(jpegQuality))
	} else {
		thumbnail := transform.Resize(i.imageData, thumbnailWidth, thumbnailHeight, transform.CatmullRom)
		err = imgio.Save(dstPath, thumbnail, imgio.JPEGEncoder(jpegQuality))
	}

	return err
}
