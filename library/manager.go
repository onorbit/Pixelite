package library

import (
	"errors"
	"sync"
)

var ErrLibraryAlreadyRegistered = errors.New("Library with given root path is already registered")
var ErrLibraryScanInProgress = errors.New("Library with given root path is being scanned")

type manager struct {
	libraries map[string]*Library
	rootPaths map[string]struct{}
	progress  map[string]struct{}
	mutex     sync.Mutex
}

var gManager manager

func (m *manager) createLibrary(rootPath string) error {
	m.mutex.Lock()

	if _, ok := m.rootPaths[rootPath]; ok == true {
		m.mutex.Unlock()
		return ErrLibraryAlreadyRegistered
	}

	if _, ok := m.progress[rootPath]; ok == true {
		m.mutex.Unlock()
		return ErrLibraryScanInProgress
	}

	m.progress[rootPath] = struct{}{}
	m.mutex.Unlock()

	newLibrary := newLibrary(rootPath)
	if err := newLibrary.scan(); err != nil {
		return err
	}

	m.mutex.Lock()
	delete(m.progress, rootPath)
	m.libraries[newLibrary.id] = &newLibrary
	m.rootPaths[rootPath] = struct{}{}
	m.mutex.Unlock()

	return nil
}

func (m *manager) getLibrary(id string) *Library {
	if _, ok := m.libraries[id]; ok == false {
		return nil
	}

	return m.libraries[id]
}

//-----------------------------------------------------------------------------
// public functions
//-----------------------------------------------------------------------------

func Initialize() error {
	gManager = manager{
		libraries: make(map[string]*Library),
		rootPaths: make(map[string]struct{}),
		progress:  make(map[string]struct{}),
		mutex:     sync.Mutex{},
	}

	return nil
}

func CreateLibrary(rootPath string) error {
	return gManager.createLibrary(rootPath)
}

func ListLibrary() []LibrarySummeryDesc {
	ret := make([]LibrarySummeryDesc, 0, len(gManager.libraries))
	for _, library := range gManager.libraries {
		summary := LibrarySummeryDesc{
			Id:   library.id,
			Desc: library.desc,
		}
		ret = append(ret, summary)
	}

	return ret
}

func GetLibrary(id string) *Library {
	return gManager.getLibrary(id)
}
