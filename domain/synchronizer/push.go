package synchronizer

import (
	"deni1688/gsync/domain"
	"log"
	"os"
)

func (g syncService) Push(sf domain.SyncFile) error {
	if sf.Name == "Gsync" {
		sf.Id = g.remoteGsyncDir
		sf.Path = g.localGsyncDir
	}

	list, err := os.ReadDir(sf.Path)
	if err != nil {
		return err
	}

	for _, file := range list {
		fullPath := domain.GetFullPath(sf.Path, file.Name())

		log.Printf("Pushing %s", fullPath)

		f := domain.SyncFile{
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
