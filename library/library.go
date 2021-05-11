package library

import (
	"fmt"
	"io/ioutil"
	"math/rand"
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

func newLibrary(rootPath string) Library {
	id := fmt.Sprintf("%08x", rand.Uint32())
	newLibrary := Library{
		id:       id,
		desc:     id,
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
		currPath := subPaths[len(subPaths)-1]
		subPaths = subPaths[0 : len(subPaths)-1]

		content, err := ioutil.ReadDir(currPath)
		if err != nil {
			return err
		}

		for _, entry := range content {
			if entry.IsDir() == true {
				path := filepath.Join(currPath, entry.Name())
				subPaths = append(subPaths, path)
			} else if image.IsImageFile(entry.Name()) == true {
				l.albums[currPath] = album.NewAlbum(currPath)
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
		desc.Albums = append(desc.Albums, album.GetPath())
	}

	return desc
}
