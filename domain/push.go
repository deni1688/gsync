package domain

import (
	"log"
	"os"
)

func (g syncService) Push(sf SyncFile) error {
	if sf.Name == "Gsync" {
		sf.Id = g.remoteGsyncDir
		sf.Path = g.localGsyncDir
	}

	list, err := os.ReadDir(sf.Path)
	if err != nil {
		return err
	}

	for _, file := range list {
		fullPath := GetFullPath(sf.Path, file.Name())

		log.Printf("Pushing %s", fullPath)

		f := SyncFile{
			Name:     file.Name(),
			Path:     fullPath,
			ParentId: sf.Id,
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
