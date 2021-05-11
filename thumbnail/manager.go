package thumbnail

import (
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/onorbit/pixelite/config"
	"github.com/onorbit/pixelite/globaldb"
)

type manager struct {
	thumbnails map[string]string
	progress   map[string]*sync.Cond
	random     *rand.Rand
	mutex      sync.Mutex
}

type doneMsg struct {
	imgPath       string
	thumbnailPath string
	err           error
}

var gManager manager

var gLetters = []rune("abcdefghijklmnopqrstuvwxyz1234567890")
var gThumbnailFileNameLen = 32

func (m *manager) getThumbnailPath(imgPath string) string {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// thumbnail already exists. return it directly.
	thumbnailPath, ok := m.thumbnails[imgPath]
	if ok == true {
		return thumbnailPath
	}

	var cond *sync.Cond
	if cond, ok = m.progress[imgPath]; ok == false {
		cond = sync.NewCond(&m.mutex)
		m.progress[imgPath] = cond

		go m.buildThumbnail(imgPath, cond)
	}

	cond.Wait()

	thumbnailPath, ok = m.thumbnails[imgPath]
	if ok == false {
		return ""
	}

	return thumbnailPath
}

func (m *manager) getRandomFileName(length int) string {
	ret := make([]rune, length)
	for i := range ret {
		ret[i] = gLetters[m.random.Intn(len(gLetters))]
	}

	return string(ret)
}

func (m *manager) buildThumbnail(imgPath string, signalCond *sync.Cond) {
	image, err := imgio.Open(imgPath)
	if err != nil {
		m.mutex.Lock()
		defer m.mutex.Unlock()

		delete(m.progress, imgPath)
		signalCond.Broadcast()

		return
	}

	bounds := image.Bounds()
	imageWidth := bounds.Dx()
	imageHeight := bounds.Dy()

	thumbnailStorePath := config.Get().Thumbnail.StorePath
	thumbnailDim := config.Get().Thumbnail.MaxDimension
	thumbnailJpegQuality := config.Get().Thumbnail.JpegQuality

	thumbnailWidth := 0
	thumbnailHeight := 0
	if imageWidth > imageHeight && imageWidth > thumbnailDim {
		thumbnailWidth = thumbnailDim
		thumbnailHeight = int(float32(imageHeight) * (float32(thumbnailDim) / float32(imageWidth)))
	} else if imageHeight > thumbnailDim {
		thumbnailWidth = int(float32(imageWidth) * (float32(thumbnailDim) / float32(imageHeight)))
		thumbnailHeight = thumbnailDim
	}

	// serve the original image directly, as it is small enough.
	if thumbnailWidth == 0 && thumbnailHeight == 0 {
		globaldb.RegisterThumbnail(imgPath, imgPath)

		m.mutex.Lock()
		defer m.mutex.Unlock()

		delete(m.progress, imgPath)
		m.thumbnails[imgPath] = imgPath
		signalCond.Broadcast()

		return
	}

	// create thumbnail image.
	thumbnailName := m.getRandomFileName(gThumbnailFileNameLen) + ".jpg"
	thumbnailPath := filepath.Join(thumbnailStorePath, thumbnailName)

	thumbnail := transform.Resize(image, thumbnailWidth, thumbnailHeight, transform.Lanczos)
	err = imgio.Save(thumbnailPath, thumbnail, imgio.JPEGEncoder(thumbnailJpegQuality))

	// register to thumbnail db once the image was made.
	if err == nil {
		globaldb.RegisterThumbnail(imgPath, thumbnailPath)
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.progress, imgPath)
	if err == nil {
		m.thumbnails[imgPath] = thumbnailPath
	}
	signalCond.Broadcast()

	return
}

func Initialize() error {
	thumbnailStorePath := config.Get().Thumbnail.StorePath
	if err := os.MkdirAll(thumbnailStorePath, 0700); err != nil {
		return err
	}

	random := rand.New(rand.NewSource(time.Hour.Nanoseconds()))
	gManager = manager{
		thumbnails: make(map[string]string),
		progress:   make(map[string]*sync.Cond),
		random:     random,
		mutex:      sync.Mutex{},
	}

	thumbnailRows, err := globaldb.LoadAllThumbnails()
	if err != nil {
		return err
	}

	for _, row := range thumbnailRows {
		gManager.thumbnails[row.ImagePath] = row.ThumbnailPath
	}

	return nil
}

func GetThumbnailPath(imgPath string) string {
	return gManager.getThumbnailPath(imgPath)
}
