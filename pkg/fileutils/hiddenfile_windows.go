// +build windows

package fileutils

import (
	"path/filepath"
	"syscall"
)

func IsHidden(filePath string) (bool, error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return false, err
	}

	utf16Path, err := syscall.UTF16PtrFromString(absPath)
	if err != nil {
		return false, err
	}

	attributes, err := syscall.GetFileAttributes(utf16Path)
	if err != nil {
		return false, err
	}

	return (attributes & syscall.FILE_ATTRIBUTE_HIDDEN) != 0, nil
}
