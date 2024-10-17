package utils

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func GetFilePath(filename string) (string, error) {
	_, currentFile, _, ok := runtime.Caller(1)
	if !ok {
		return "", fmt.Errorf("unable to get caller information")
	}

	dir := filepath.Dir(currentFile)

	fullPath := filepath.Join(dir, filename)

	return fullPath, nil
}
