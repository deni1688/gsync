package domain

import (
	"fmt"
	"log"
	"os"
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

func (g gsyncService) Pull(path string) error {
	info := FileInfo{}

	if path == "" {
		info.Id = g.remoteGsyncDir
	} else {
		info.Name = path
		inf, err := g.store.GetFile(info)
		if err != nil {
			return err
		}
		info = inf
	}

	files, err := g.store.ListFiles(info.Id)
	if err != nil {
		return err
	}

	for _, file := range files {
		fullPath := fmt.Sprintf("%s/%s", path, file.Name)

		if g.store.IsDir(file) {
			if path != "Gsync" {
				err = createDir(fullPath)
				if err != nil {
					return err
				}
			}

			return g.Pull(info.Id)
		}

		info, err := g.store.GetFile(file)
		if err != nil {
			return err
		}

		return os.WriteFile(fullPath, info.Data, 0700)
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
