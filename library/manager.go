package library

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"

	"github.com/onorbit/pixelite/database/globaldb"
)

var ErrLibraryAlreadyRegistered = errors.New("Library with given root path is already registered")
var ErrLibraryScanInProgress = errors.New("Library with given root path is being scanned")
var ErrLibraryNotFound = errors.New("Library with given ID is not found")

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

	id := fmt.Sprintf("%08x", rand.Uint32())
	// TODO : receive desc from user
	newLibrary := newLibrary(id, rootPath, id)
	if err := newLibrary.scan(); err != nil {
		return err
	}

	m.mutex.Lock()
	delete(m.progress, rootPath)
	m.addLibrary(&newLibrary)
	m.mutex.Unlock()

	// TODO : what to do if some error happens in here?
	err := globaldb.InsertLibrary(id, rootPath, id)
	return err
}

func (m *manager) addLibrary(library *Library) {
	m.libraries[library.id] = library
	m.rootPaths[library.rootPath] = struct{}{}
}

func (m *manager) getLibrary(id string) *Library {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.libraries[id]; ok == false {
		return nil
	}

	return m.libraries[id]
}

func (m *manager) deleteLibrary(id string) error {
	m.mutex.Lock()

	if _, ok := m.libraries[id]; ok == true {
		delete(m.libraries, id)
	} else {
		m.mutex.Unlock()
		return ErrLibraryNotFound
	}

	m.mutex.Unlock()

	// TODO : what to do if some error happens in here?
	err := globaldb.DeleteLibrary(id)
	return err
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

	libraryRows, err := globaldb.LoadAllLibraries()
	if err != nil {
		return err
	}

	for _, row := range libraryRows {
		library := newLibrary(row.ID, row.RootPath, row.Desc)

		// TODO : library should be re-scanned as albums are not saved in its own db, by now.
		if err := library.scan(); err != nil {
			return err
		}

		gManager.addLibrary(&library)
	}

	return nil
}

func CreateLibrary(rootPath string) error {
	return gManager.createLibrary(rootPath)
}

func GetLibrary(id string) *Library {
	return gManager.getLibrary(id)
}

func DeleteLibrary(id string) error {
	return gManager.deleteLibrary(id)
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
