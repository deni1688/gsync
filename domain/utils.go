package domain

import (
	"os"
	"strings"
)

func SyncFileListContains(list []SyncFile, name string) bool {
	for _, file := range list {
		if file.Name == name {
			return true
		}
	}

	return false
}

func GetFullPath(items ...string) string {
	return strings.Join(items, "/")
}

func CreateDir(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, 0700)
	}

	return err
}
