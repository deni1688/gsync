package synchronizer

import (
	"deni1688/gsync/domain"
	"log"
	"os"
)

func (g syncService) Pull(sf domain.SyncFile) error {
	if sf.Name == "Gsync" {
		sf.Id = g.remoteGsyncDir
		sf.Path = g.localGsyncDir
	}

	files, err := g.drive.ListFiles(sf)
	if err != nil {
		return err
	}

	if err = g.addFilesFromRemote(sf, files); err != nil {
		return err
	}

	return g.removeFilesFromLocal(sf, files)
}

func (g syncService) removeFilesFromLocal(sf domain.SyncFile, files []domain.SyncFile) error {
	list, err := os.ReadDir(sf.Path)
	if err != nil {
		return err
	}

	for _, file := range list {
		name := file.Name()

		if domain.SyncFileListContains(files, name) {
			continue
		}

		fullPath := domain.GetFullPath(sf.Path, name)

		log.Printf("Removing %s", fullPath)

		if file.IsDir() {
			err = os.RemoveAll(fullPath)
		} else {
			err = os.Remove(fullPath)
		}
	}

	return err
}

func (g syncService) addFilesFromRemote(sf domain.SyncFile, files []domain.SyncFile) error {
	var err error

	for _, file := range files {
		fullPath := domain.GetFullPath(sf.Path, file.Name)
		log.Printf("Pulling %s", fullPath)

		if g.drive.IsDir(file) {
			err = domain.CreateDir(fullPath)
			if err != nil {
				return err
			}

			file.Path = fullPath

			if err = g.Pull(file); err != nil {
				return err
			}

			continue
		}

		file.Data, err = g.drive.GetFile(file)
		if err != nil {
			return err
		}

		err = os.WriteFile(fullPath, file.Data, 0700)
	}

	return err
}
