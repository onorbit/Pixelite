package library

import (
	"errors"
	"os"
	"path"
	"sync"

	"github.com/onorbit/pixelite/database/globaldb"
	"github.com/onorbit/pixelite/database/librarydb"
	"github.com/onorbit/pixelite/pkg/log"
)

var ErrLibraryAlreadyMounted = errors.New("Library with given root path is already mounted")
var ErrLibraryScanInProgress = errors.New("Library with given root path is being scanned")
var ErrLibraryDBFileNotFound = errors.New("Library database file not found")
var ErrLibraryNotFound = errors.New("Library with given ID is not found")

type manager struct {
	libraries map[string]*Library
	mutex     sync.Mutex
}

var gManager manager

func (m *manager) mountLibrary(rootPath string, isNewLibrary bool) error {
	// check if the database file exists.
	libDBPath := path.Join(rootPath, "library.sqlite3")
	if stat, err := os.Stat(libDBPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if !isNewLibrary {
				return ErrLibraryDBFileNotFound
			}
		} else {
			return err
		}
	} else if stat.IsDir() {
		return ErrLibraryDBFileNotFound
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// check if the path is already mounted.
	for _, lib := range m.libraries {
		if lib.rootPath == rootPath {
			return ErrLibraryAlreadyMounted
		}
	}

	// load or initialize libraryDB.
	libDB, err := librarydb.LoadLibraryDB(libDBPath)
	if err != nil {
		return err
	}

	libraryID := libDB.GetLibraryID()
	libraryDesc, _ := libDB.GetMetadata(librarydb.MetadataKeyLibraryTitle)
	if len(libraryDesc) == 0 {
		libraryDesc = rootPath
	}

	newLibrary := newLibrary(libraryID, rootPath, libraryDesc)
	if err := newLibrary.scan(); err != nil {
		return err
	}
	m.libraries[newLibrary.id] = newLibrary
	if isNewLibrary {
		err = globaldb.InsertLibrary(rootPath)
		// TODO : what to do if some error happens in here?
		if err != nil {
			return err
		}
	}

	log.Info("library [%s] with root path [%s] is successfully mounted", libraryID, rootPath)
	return nil
}

func (m *manager) getLibrary(id string) *Library {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.libraries[id]; !ok {
		return nil
	}

	return m.libraries[id]
}

func (m *manager) unmountLibrary(id string) error {
	m.mutex.Lock()

	library, ok := m.libraries[id]
	if ok {
		delete(m.libraries, id)
	} else {
		m.mutex.Unlock()
		log.Warn("failed to unmount library [%s] as not found", id)
		return ErrLibraryNotFound
	}

	librarydb.UnloadLibraryDB(id)
	m.mutex.Unlock()

	err := globaldb.DeleteLibrary(library.rootPath)
	if err != nil {
		// TODO : what to do if some error happens in here?
	}

	log.Info("library [%s] unmounted successfully", id)

	return err
}

//-----------------------------------------------------------------------------
// public functions
//-----------------------------------------------------------------------------

func Initialize() error {
	gManager = manager{
		libraries: make(map[string]*Library),
		mutex:     sync.Mutex{},
	}

	libraryRows, err := globaldb.LoadAllLibraries()
	if err != nil {
		return err
	}

	for _, row := range libraryRows {
		if err = gManager.mountLibrary(row.RootPath, false); err != nil {
			if err == ErrLibraryDBFileNotFound {
				log.Warn("library DB file not found at [%s], unmounting", row.RootPath)
				globaldb.DeleteLibrary(row.RootPath)

				continue
			}
			return err
		}
	}

	return nil
}

func MountLibrary(rootPath string) error {
	return gManager.mountLibrary(rootPath, true)
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
			Id:    library.id,
			Title: library.title,
		}
		ret = append(ret, summary)
	}

	return ret
}
