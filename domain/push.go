package domain

import (
	"fmt"
	"log"
	"os"
)

func (g syncService) Push(dir SyncFile) error {
	if dir.Name == "Gsync" {
		dir.Id = g.remoteGsyncDir
		dir.Path = g.localGsyncDir
	}

	if dir.Id == "" {
		return fmt.Errorf("dir id is required")
	}

	list, err := os.ReadDir(dir.Path)
	if err != nil {
		return err
	}

	for _, file := range list {
		fullPath := GetPathFrom(dir.Path, file.Name())

		log.Printf("Pushing %s", fullPath)

		f := SyncFile{
			Name:     file.Name(),
			Path:     fullPath,
			ParentId: dir.Id,
		}

		if file.IsDir() {
			f, err = g.drive.CreateDir(f)
			if err != nil {
				return err
			}

			if err = g.Push(f); err != nil {
				return err
			}

			continue
		}

		f.Data, err = os.ReadFile(fullPath)
		if err != nil {
			return err
		}

		f, err = g.drive.CreateFile(f)
	}

	// TODO: Remove files from remote that are not in local

	return nil
}
