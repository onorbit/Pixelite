package library

import (
	"io/ioutil"
	"path/filepath"

	"github.com/onorbit/pixelite/album"
	"github.com/onorbit/pixelite/image"
)

type Library struct {
	id       string
	desc     string
	rootPath string
	albums   map[string]album.Album
}

func newLibrary(id, rootPath, desc string) Library {
	newLibrary := Library{
		id:       id,
		desc:     desc,
		rootPath: rootPath,
		albums:   make(map[string]album.Album),
	}

	return newLibrary
}

func (l Library) scan() error {
	if len(l.albums) != 0 {
		l.albums = make(map[string]album.Album)
	}

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

		for _, entry := range content {
			if entry.IsDir() == true {
				// found a directory. push to subpath list for further traverse.
				path := filepath.Join(currPath, entry.Name())
				subPaths = append(subPaths, path)
			} else if image.IsImageFile(entry.Name()) == true {
				// found an image. register this path as an Album.
				albumID, _ := filepath.Rel(l.rootPath, currPath)
				albumID = filepath.ToSlash(albumID)
				l.albums[albumID] = album.NewAlbum(albumID, currPath)
				break
			}
		}
	}

	return nil
}

//-----------------------------------------------------------------------------
// public functions
//-----------------------------------------------------------------------------

func (l Library) Describe() LibraryDesc {
	desc := LibraryDesc{
		Id:     l.id,
		Desc:   l.desc,
		Albums: make([]string, 0, len(l.albums)),
	}

	for _, album := range l.albums {
		desc.Albums = append(desc.Albums, album.GetID())
	}

	return desc
}

// TODO : returning pointer here could be dangerous. need to fix.
func (l Library) GetAlbum(albumID string) *album.Album {
	if ret, ok := l.albums[albumID]; ok {
		return &ret
	}

	return nil
}
