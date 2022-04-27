package domain

import (
	"log"
	"os"
)

func (g gsyncService) Push(fi FileInfo) error {
	if fi.Name == "Gsync" {
		fi.Id = g.remoteGsyncDir
		fi.Path = g.localGsyncDir
	}

	list, err := os.ReadDir(fi.Path)
	if err != nil {
		return err
	}

	for _, file := range list {
		fullPath := getFullPath(fi.Path, file.Name())

		log.Printf("Pushing %s", fullPath)

		f := FileInfo{
			Name:     file.Name(),
			Path:     fullPath,
			ParentId: fi.Id,
		}

		if file.IsDir() {
			f, err = g.store.CreateDir(f)
			if err != nil {
				return err
			}

			err = g.Push(f)
			if err != nil {
				return err
			}

			continue
		}

		f.Data, err = os.ReadFile(fullPath)
		if err != nil {
			return err
		}

		f, err = g.store.CreateFile(f)
		if err != nil {
			return err
		}
	}

	// TODO: Remove files from remote that are not in local

	return nil
}
