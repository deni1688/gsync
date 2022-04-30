package syncer

import (
	"os"
	"strings"
)

func FileListContains(fileList []SyncFile, name string) bool {
	for _, file := range fileList {
		if file.Name == name {
			return true
		}
	}

	return false
}

func GetPathFrom(pathItems ...string) string {
	return strings.Join(pathItems, "/")
}

func CreateDir(dirFullPath string) error {
	_, err := os.Stat(dirFullPath)
	if os.IsNotExist(err) {
		err = os.Mkdir(dirFullPath, 0700)
	}

	return err
}
