package domain

import (
	"fmt"
	"log"
	"os"
)

type gsyncService struct {
	localGsyncDir string
	store         SynchronizableStoreContract
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

	return &gsyncService{localGsyncDir, store}
}

func (g gsyncService) Pull(path string) error {
	if path == "" {
		path = g.localGsyncDir
	}

	files, err := g.store.ListFiles(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		fullPath := fmt.Sprintf("%s/%s", path, file.Name)

		if g.store.IsDir(file) {
			err = createDir(fullPath)
			if err != nil {
				return err
			}

			return g.Pull(fullPath)
		}

		data, err := g.store.GetFile(file)
		if err != nil {
			return err
		}

		return os.WriteFile(fullPath, data, 0700)
	}

	return nil
}

func createDir(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, 0700)
	}

	return err
}

func (g gsyncService) Push(option ...SyncOption) error {
	//TODO implement me
	panic("implement me")
}

func (g gsyncService) Sync(option ...SyncOption) error {
	//TODO implement me
	panic("implement me")
}
