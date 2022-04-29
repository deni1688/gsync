package syncer

import (
	"log"
	"os"
)

func (g syncService) Pull(sf SyncFile) error {
	if sf.Name == "Gsync" {
		sf.Id = g.remoteGsyncDir
		sf.Path = g.localGsyncDir
	}

	files, err := g.drive.ListFiles(sf)
	if err != nil {
		return err
	}

	return g.addFilesFromRemote(sf, files)
}

func (g syncService) removeFilesFromLocal(sf SyncFile, files []SyncFile) error {
	list, err := os.ReadDir(sf.Path)
	if err != nil {
		return err
	}

	for _, file := range list {
		name := file.Name()

		if SyncFileListContains(files, name) {
			continue
		}

		fullPath := GetFullPath(sf.Path, name)

		log.Printf("Removing %s", fullPath)

		if file.IsDir() {
			err = os.RemoveAll(fullPath)
		} else {
			err = os.Remove(fullPath)
		}
	}

	return err
}

func (g syncService) addFilesFromRemote(sf SyncFile, files []SyncFile) error {
	errCh := make(chan error, 1)

	for _, file := range files {
		go func(errCh chan error, file SyncFile) {
			fullPath := GetFullPath(sf.Path, file.Name)
			log.Printf("Pulling %s", fullPath)

			if g.drive.IsDir(file) {
				err := CreateDir(fullPath)
				if err != nil {
					errCh <- err
				}
				file.Path = fullPath
				err = g.Pull(file)
				if err != nil {
					errCh <- err
				}
			} else {
				data, err := g.drive.GetFile(file)
				if err != nil {
					errCh <- err
				}

				file.Data = data

				if err = os.WriteFile(fullPath, file.Data, 0700); err != nil {
					errCh <- err
				}
			}
			errCh <- nil
		}(errCh, file)
	}

	return <-errCh
}
