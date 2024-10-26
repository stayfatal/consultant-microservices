package utils

import (
	"os"
	"strings"
)

func GetPath(relativePath string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dirs := strings.Split(wd, "/")
	var path string
	for i := 1; i < len(dirs); i++ {
		path += "/" + dirs[i]
		if dirs[i] == "consultant-microservices" {
			break
		}
	}

	return path + "/" + relativePath, nil
}
