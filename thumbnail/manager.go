package thumbnail

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/onorbit/pixelite/config"
	"github.com/onorbit/pixelite/database/globaldb"
	"github.com/onorbit/pixelite/image"
)

type manager struct {
	thumbnails map[string]string
	progress   map[string]*sync.Cond
	random     *rand.Rand
	mutex      sync.Mutex
}

var gManager manager

var gLetters = []rune("abcdefghijklmnopqrstuvwxyz1234567890")
var gThumbnailFileNameLen = 32

func (m *manager) getThumbnailPath(fileName, albumPath, albumID string) string {
	// make thumbnail directory.
	albumIDHashArr := md5.Sum([]byte(albumID))
	albumIDHash := hex.EncodeToString(albumIDHashArr[:])
	thumbnailDir := filepath.Join(config.Get().Thumbnail.StorePath, albumIDHash)

	if err := os.MkdirAll(thumbnailDir, 0700); err != nil {
		// TODO : handle the error properly.
		return ""
	}

	// make thumbnail file path.
	thumbnailName := strings.TrimSuffix(fileName, filepath.Ext(fileName)) + ".jpg"
	thumbnailPath := filepath.Join(thumbnailDir, thumbnailName)

	// make original image file path.
	imgPath := path.Join(albumPath, fileName)

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// thumbnail already exists. return it directly.
	existThumbnailPath, ok := m.thumbnails[imgPath]
	if ok {
		return existThumbnailPath
	}

	var cond *sync.Cond
	if cond, ok = m.progress[imgPath]; !ok {
		cond = sync.NewCond(&m.mutex)
		m.progress[imgPath] = cond

		go m.buildThumbnail(imgPath, thumbnailPath, cond)
	}

	cond.Wait()

	thumbnailPath, ok = m.thumbnails[imgPath]
	if !ok {
		return ""
	}

	return thumbnailPath
}

func (m *manager) buildThumbnail(imgPath, thumbnailPath string, signalCond *sync.Cond) {
	// get parameters for making thumbnail.
	thumbnailDim := config.Get().Thumbnail.MaxDimension
	thumbnailJpegQuality := config.Get().Thumbnail.JpegQuality

	// make actual thumbnail.
	err := image.MakeThumbnail(imgPath, thumbnailPath, thumbnailDim, thumbnailJpegQuality)

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
}

//-----------------------------------------------------------------------------
// public functions
//-----------------------------------------------------------------------------

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

func GetThumbnailPath(fileName, albumPath, albumID string) string {
	return gManager.getThumbnailPath(fileName, albumPath, albumID)
}
