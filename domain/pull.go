package domain

import (
	"log"
	"os"
)

func (g gsyncService) Pull(fi FileInfo) error {
	if fi.Name == "Gsync" {
		fi.Id = g.remoteGsyncDir
		fi.Path = g.localGsyncDir
	}

	files, err := g.store.ListFiles(fi)
	if err != nil {
		return err
	}

	if err = g.addFilesFromRemote(fi, files); err != nil {
		return err
	}

	return g.removeFilesFromLocal(fi, files)
}

func (g gsyncService) removeFilesFromLocal(fi FileInfo, files []FileInfo) error {
	list, err := os.ReadDir(fi.Path)
	if err != nil {
		return err
	}

	for _, file := range list {
		name := file.Name()

		if fileInfoListContains(files, name) {
			continue
		}

		fullPath := getFullPath(fi.Path, name)

		log.Printf("Removing %s", fullPath)

		if file.IsDir() {
			err = os.RemoveAll(fullPath)
		} else {
			err = os.Remove(fullPath)
		}
	}

	return err
}

func (g gsyncService) addFilesFromRemote(fi FileInfo, files []FileInfo) error {
	var err error

	for _, file := range files {
		fullPath := getFullPath(fi.Path, file.Name)
		log.Printf("Pulling %s", fullPath)

		if g.store.IsDir(file) {
			err = createDir(fullPath)
			if err != nil {
				return err
			}

			file.Path = fullPath

			if err = g.Pull(file); err != nil {
				return err
			}

			continue
		}

		file.Data, err = g.store.GetFile(file)
		if err != nil {
			return err
		}

		err = os.WriteFile(fullPath, file.Data, 0700)
		if err != nil {
			return err
		}
	}

	return nil
}
