package domain

import (
	"log"
	"os"
	"strings"
)

type gsyncService struct {
	remoteGsyncDir string
	localGsyncDir  string
	store          SynchronizableStoreContract
}

func NewGsyncService(localGsyncDir string, store SynchronizableStoreContract) GsyncServiceContract {
	if localGsyncDir == "" {
		localGsyncDir = os.Getenv("HOME") + "/Gsync"
	}

	if _, err := os.Stat(localGsyncDir); os.IsNotExist(err) {
		log.Println("Creating local gsync directory...")

		err = createDir(localGsyncDir)
		if err != nil {
			log.Fatalf("Error creating local gsync directory: %v", err)
		}

		log.Println("Local gsync directory created in ", localGsyncDir)
	}

	info, err := store.CreateDirectory("Gsync")
	if err != nil {
		log.Fatalf("Error creating remote Gsync directory: %v", err)
	}

	return &gsyncService{info.Id, localGsyncDir, store}
}

func (g gsyncService) Pull(fi FileInfo) error {
	if fi.Name == "Gsync" {
		fi.Id = g.remoteGsyncDir
		fi.Path = g.localGsyncDir
	}

	files, err := g.store.ListFiles(fi.Path, fi.Id)
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
		if err = os.Remove(getFullPath(fi.Path, name)); err != nil {
			return err
		}
	}

	return nil
}

func (g gsyncService) addFilesFromRemote(fi FileInfo, files []FileInfo) error {
	var err error

	for _, file := range files {
		fullPath := getFullPath(fi.Path, file.Name)

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

func fileInfoListContains(list []FileInfo, name string) bool {
	for _, file := range list {
		if file.Name == name {
			return true
		}
	}

	return false
}

func getFullPath(items ...string) string {
	return strings.Join(items, "/")
}

func createDir(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, 0700)
	}

	return err
}

func (g gsyncService) Push(info FileInfo) error {
	//TODO implement me
	panic("implement me")
}

func (g gsyncService) Sync(info FileInfo) error {
	var err error

	err = g.Push(info)
	err = g.Pull(info)

	return err
}
