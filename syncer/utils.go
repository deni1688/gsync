package syncer

import (
	"os"
	"strings"
)

func fileListContains(fileList []SyncFile, name string) bool {
	for _, file := range fileList {
		if file.Name == name {
			return true
		}
	}

	return false
}

func getPath(pathItems ...string) string {
	return strings.Join(pathItems, "/")
}

func createDirs(paths ...string) error {
	for _, path := range paths {
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			err = os.Mkdir(path, 0700)
		} else {
			return err
		}
	}

	return nil
}
