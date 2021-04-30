package library

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"sync"

	"azurestud.io/pixelite/image"
)

var ErrLibraryAlreadyRegistered = errors.New("Library with given root path is already registered")
var ErrLibraryScanInProgress = errors.New("Library with given root path is being scanned")

type manager struct {
	libraries map[string]Library
	progress  map[string]struct{}
	mutex     sync.Mutex
}

var gManager manager

func (m *manager) createLibrary(rootPath string) error {
	m.mutex.Lock()

	if _, ok := m.libraries[rootPath]; ok == true {
		m.mutex.Unlock()
		return ErrLibraryAlreadyRegistered
	}

	if _, ok := m.progress[rootPath]; ok == true {
		m.mutex.Unlock()
		return ErrLibraryScanInProgress
	}

	m.progress[rootPath] = struct{}{}
	m.mutex.Unlock()

	newLibrary := Library{
		RootPath: rootPath,
		Albums:   make(map[string]struct{}),
	}

	// perform scan.
	dirs := make([]string, 0, 1)
	dirs = append(dirs, rootPath)

	for len(dirs) != 0 {
		currPath := dirs[len(dirs)-1]
		dirs = dirs[0 : len(dirs)-1]

		content, err := ioutil.ReadDir(currPath)
		if err != nil {
			return err
		}

		for _, entry := range content {
			if entry.IsDir() == true {
				path := filepath.Join(currPath, entry.Name())
				dirs = append(dirs, path)
			} else if image.IsImageFile(entry.Name()) == true {
				newLibrary.Albums[currPath] = struct{}{}
				break
			}
		}
	}

	m.mutex.Lock()
	delete(m.progress, rootPath)
	m.libraries[rootPath] = newLibrary
	m.mutex.Unlock()

	return nil
}

func Initialize() error {
	gManager = manager{
		libraries: make(map[string]Library),
		progress:  make(map[string]struct{}),
		mutex:     sync.Mutex{},
	}

	return nil
}

func CreateLibrary(rootPath string) error {
	return gManager.createLibrary(rootPath)
}
