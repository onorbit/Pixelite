package library

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/onorbit/pixelite/album"
	"github.com/onorbit/pixelite/database/librarydb"
	"github.com/onorbit/pixelite/image"
	"github.com/onorbit/pixelite/pkg/fileutils"
	"github.com/onorbit/pixelite/pkg/log"
)

var ErrLibraryDBNotFound = errors.New("librarydb for the library not found")

type Library struct {
	id       string
	title    string
	rootPath string
	albums   map[string]album.Album
	mutex    sync.Mutex
}

func newLibrary(id, rootPath, title string) *Library {
	newLibrary := &Library{
		id:       id,
		title:    title,
		rootPath: rootPath,
		albums:   make(map[string]album.Album),
	}

	return newLibrary
}

func (l *Library) scan() error {
	newAlbums := make(map[string]album.Album)
	subPaths := make([]string, 0, 1)
	subPaths = append(subPaths, l.rootPath)

	for len(subPaths) != 0 {
		// pop an entry from subpath list.
		currPath := subPaths[len(subPaths)-1]
		subPaths = subPaths[0 : len(subPaths)-1]

		content, err := ioutil.ReadDir(currPath)
		if err != nil {
			return err
		}

		isRegistered := false
		for _, entry := range content {
			if entry.IsDir() == true {
				// found a directory. push to subpath list for further traverse.
				path := filepath.Join(currPath, entry.Name())
				subPaths = append(subPaths, path)
			} else if image.IsImageFile(entry.Name()) == true {
				isHidden, err := fileutils.IsHidden(entry.Name())
				if isHidden || err != nil {
					continue
				}

				// found an image. register this path as an Album.
				if !isRegistered {
					albumID, _ := filepath.Rel(l.rootPath, currPath)
					albumID = filepath.ToSlash(albumID)
					newAlbums[albumID] = album.NewAlbum(albumID, currPath, entry.Name())

					isRegistered = true
				}
			}
		}
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.albums = newAlbums

	return nil
}

//-----------------------------------------------------------------------------
// public functions
//-----------------------------------------------------------------------------

func (l *Library) Describe() LibraryDesc {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	desc := LibraryDesc{
		Id:     l.id,
		Title:  l.title,
		Albums: make([]string, 0, len(l.albums)),
	}

	for _, album := range l.albums {
		desc.Albums = append(desc.Albums, album.GetID())
	}

	return desc
}

// TODO : returning pointer here could be dangerous. need to fix.
func (l *Library) GetAlbum(albumID string) *album.Album {
	if ret, ok := l.albums[albumID]; ok {
		return &ret
	}

	return nil
}

func (l *Library) Rescan() error {
	return l.scan()
}

func (l *Library) SetTitle(title string) error {
	l.title = title

	libdb := librarydb.GetLibraryDB(l.id)
	if libdb == nil {
		log.Error("failed to find libraryDB for library [%s] while changing title", l.id)
		return ErrLibraryDBNotFound
	}

	err := libdb.SetMetadata(librarydb.MetadataKeyLibraryTitle, title)
	if err != nil {
		log.Error("failed to set title [%s] to librarydb [%s] - error [%v]", title, l.id, err.Error())
		return err
	}

	return nil
}
