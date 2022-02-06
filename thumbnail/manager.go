package thumbnail

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"os"
	"path/filepath"
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

func (m *manager) getThumbnailPath(imgPath, albumID string) string {
	albumIDHashArr := md5.Sum([]byte(albumID))
	albumIDHash := hex.EncodeToString(albumIDHashArr[:])

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// thumbnail already exists. return it directly.
	thumbnailPath, ok := m.thumbnails[imgPath]
	if ok {
		return thumbnailPath
	}

	var cond *sync.Cond
	if cond, ok = m.progress[imgPath]; !ok {
		cond = sync.NewCond(&m.mutex)
		m.progress[imgPath] = cond

		go m.buildThumbnail(imgPath, albumIDHash, cond)
	}

	cond.Wait()

	thumbnailPath, ok = m.thumbnails[imgPath]
	if !ok {
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

func (m *manager) buildThumbnail(imgPath, albumIDHash string, signalCond *sync.Cond) {
	// make path for thumbnail.
	thumbnailStorePath := config.Get().Thumbnail.StorePath
	thumbnailPath := filepath.Join(thumbnailStorePath, albumIDHash)

	if err := os.MkdirAll(thumbnailPath, 0700); err != nil {
		// TODO : handle the error properly.
		return
	}

	thumbnailName := m.getRandomFileName(gThumbnailFileNameLen)
	outputPath := filepath.Join(thumbnailPath, thumbnailName) + ".jpg"

	// get parameters for making thumbnail.
	thumbnailDim := config.Get().Thumbnail.MaxDimension
	thumbnailJpegQuality := config.Get().Thumbnail.JpegQuality

	// make actual thumbnail.
	err := image.MakeThumbnail(imgPath, outputPath, thumbnailDim, thumbnailJpegQuality)

	if err == nil {
		globaldb.RegisterThumbnail(imgPath, outputPath)
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.progress, imgPath)
	if err == nil {
		m.thumbnails[imgPath] = outputPath
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

func GetThumbnailPath(imgPath, albumID string) string {
	return gManager.getThumbnailPath(imgPath, albumID)
}
