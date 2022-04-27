package domain

import (
	"os"
	"strings"
)

func fileInfoListContains(list []FileInfo, name string) bool {
	for _, file := range list {
		if file.Name == name {
			return true
		}
	}

	return false
}

func getFullPath(items ...string) string {
	return strings.Join(items, "/")
}

func createDir(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, 0700)
	}

	return err
}
