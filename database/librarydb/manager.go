package librarydb

import (
	"errors"
	"sync"
)

var (
	ErrNotFound      = errors.New("libraryDB not found")
	ErrAlreadyLoaded = errors.New("libraryDB already loaded")
)

type manager struct {
	libraryDBs map[string]*LibraryDB
	mutex      sync.Mutex
}

var gManager *manager

func (m *manager) LoadLibraryDB(libraryID, dbFilePath string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.libraryDBs[libraryID]; ok {
		return ErrAlreadyLoaded
	}

	// TODO : implement this.
	return nil
}

func (m *manager) UnloadLibraryDB(libraryID string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.libraryDBs[libraryID]; !ok {
		return ErrNotFound
	}

	delete(m.libraryDBs, libraryID)
	return nil
}

func (m *manager) GetLibraryDB(libraryID string) *LibraryDB {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	db := m.libraryDBs[libraryID]
	return db
}

//-----------------------------------------------------------------------------
// public functions
//-----------------------------------------------------------------------------

func Initialize() error {
	gManager = &manager{
		libraryDBs: make(map[string]*LibraryDB),
		mutex:      sync.Mutex{},
	}

	return nil
}

func LoadLibraryDB(libraryID, dbFilePath string) error {
	return gManager.LoadLibraryDB(libraryID, dbFilePath)
}

func UnloadLibraryDB(libraryID string) error {
	return gManager.UnloadLibraryDB(libraryID)
}

func GetLibraryDB(libraryID string) *LibraryDB {
	return gManager.GetLibraryDB(libraryID)
}
