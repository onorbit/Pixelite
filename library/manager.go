package library

import (
	"errors"
	"path"
	"sync"

	"github.com/onorbit/pixelite/database/globaldb"
	"github.com/onorbit/pixelite/database/librarydb"
	"github.com/onorbit/pixelite/pkg/log"
)

var ErrLibraryAlreadyMounted = errors.New("Library with given root path is already mounted")
var ErrLibraryScanInProgress = errors.New("Library with given root path is being scanned")
var ErrLibraryNotFound = errors.New("Library with given ID is not found")

type manager struct {
	libraries map[string]*Library
	rootPaths map[string]struct{}
	progress  map[string]struct{}
	mutex     sync.Mutex
}

var gManager manager

func (m *manager) mountLibrary(rootPath string) error {
	m.mutex.Lock()

	// check if the path is already mounted.
	if _, ok := m.rootPaths[rootPath]; ok == true {
		m.mutex.Unlock()
		return ErrLibraryAlreadyMounted
	}

	// or being mounted and scanned.
	if _, ok := m.progress[rootPath]; ok == true {
		m.mutex.Unlock()
		return ErrLibraryScanInProgress
	}

	// load or initialize libraryDB.
	libDBPath := path.Join(rootPath, "library.sqlite3")
	id, err := librarydb.LoadLibraryDB(libDBPath)
	if err != nil {
		m.mutex.Unlock()
		return err
	}

	m.progress[rootPath] = struct{}{}
	m.mutex.Unlock()

	// TODO : receive desc from user
	newLibrary := newLibrary(id, rootPath, id)
	if err := newLibrary.scan(); err != nil {
		return err
	}

	m.mutex.Lock()
	delete(m.progress, rootPath)
	m.addLibrary(newLibrary)
	m.mutex.Unlock()

	// TODO : what to do if some error happens in here?
	err = globaldb.InsertLibrary(id, rootPath, id)
	if err != nil {
		return err
	}

	log.Info("library [%s] with root path [%s] is successfully mounted", id, rootPath)
	return nil
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

func (m *manager) unmountLibrary(id string) error {
	m.mutex.Lock()

	if _, ok := m.libraries[id]; ok == true {
		delete(m.libraries, id)
	} else {
		m.mutex.Unlock()
		return ErrLibraryNotFound
	}

	librarydb.UnloadLibraryDB(id)

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

		gManager.addLibrary(library)
	}

	return nil
}

func MountLibrary(rootPath string) error {
	return gManager.mountLibrary(rootPath)
}

func GetLibrary(id string) *Library {
	return gManager.getLibrary(id)
}

func UnmountLibrary(id string) error {
	return gManager.unmountLibrary(id)
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
