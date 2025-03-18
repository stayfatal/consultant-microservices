package utils

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func GetPath(relativePath string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if filepath.Base(wd) == "consultant-microservices" {
			return filepath.Join(wd, relativePath), nil
		}

		parent := filepath.Dir(wd)
		if parent == wd { // Достигли корня файловой системы
			return "", errors.New("consultant-microservices directory not found")
		}
		wd = parent
	}
}
