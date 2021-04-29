package image

import (
	"path/filepath"
	"strings"
)

var imageExt = []string{
	".jpg",
	".png",
}

func IsImageFile(fileName string) bool {
	fileExt := filepath.Ext(fileName)
	fileExt = strings.ToLower(fileExt)
	for _, ext := range imageExt {
		if ext == fileExt {
			return true
		}
	}

	return false
}
