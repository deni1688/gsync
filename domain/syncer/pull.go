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

	if err = g.removeFilesFromLocal(sf, files); err != nil {
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
		fullPath := GetFullPath(sf.Path, name)

		if SyncFileListContains(files, name) {
			continue
		}

		log.Printf("Removing %s", fullPath)

		if file.IsDir() {
			if err = os.RemoveAll(fullPath); err != nil {
				return err
			}
		} else {
			if err = os.Remove(fullPath); err != nil {
				return err
			}
		}
	}

	return err
}

func (g syncService) addFilesFromRemote(sf SyncFile, files []SyncFile) error {
	for _, file := range files {
		fullPath := GetFullPath(sf.Path, file.Name)

		log.Printf("Pulling %s", fullPath)

		if g.drive.IsDir(file) {
			if err := CreateDir(fullPath); err != nil {
				return err
			}

			file.Path = fullPath

			err := g.Pull(file)
			if err != nil {
				return err
			}

			continue
		}

		data, err := g.drive.GetFile(file)
		if err != nil {
			return err
		}

		return os.WriteFile(fullPath, data, 0700)
	}

	return nil
}
