package media

import "errors"

func movieFileLoader(srcPath string) (MediaFile, error) {
	ret := &movieFile{}
	return ret, nil
}

func registerMovieLoaders() {
	gMediaFileLoaders[".mp4"] = movieFileLoader
}

type movieFile struct {
}

func (m *movieFile) MakeThumbnail(dstPath string, thumbnailSize, jpegQuality int, squareCrop bool) error {
	return errors.New("not implemented")
}
