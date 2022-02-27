// +build !windows

package fileutils

import "path/filepath"

func IsHidden(filePath string) (bool, error) {
	fileName := filepath.Base(filePath)
	return fileName[0] == '.', nil
}
