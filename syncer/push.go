package syncer

import (
	"fmt"
	"log"
	"os"
)

func (gs gsyncService) Push(dir SyncFile) error {
	if dir.Name == "Gsync" {
		dir.Id = gs.remoteGsyncDir
		dir.Path = gs.localGsyncDir
	}

	if dir.Id == "" {
		return fmt.Errorf("dir id is required")
	}

	list, err := os.ReadDir(dir.Path)
	if err != nil {
		return err
	}

	for _, file := range list {
		fullPath := getPath(dir.Path, file.Name())

		log.Printf("Pushing %s", fullPath)

		f := SyncFile{
			Name:     file.Name(),
			Path:     fullPath,
			ParentId: dir.Id,
		}

		if file.IsDir() {
			f, err = gs.syncProvider.CreateDir(f)
			if err != nil {
				return err
			}

			if err = gs.Push(f); err != nil {
				return err
			}

			continue
		}

		f.Data, err = os.ReadFile(fullPath)
		if err != nil {
			return err
		}

		f, err = gs.syncProvider.CreateFile(f)
	}

	// TODO: Remove files from remote that are not in local

	return nil
}
