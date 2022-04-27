package domain

import (
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

func (g gsyncService) Pull(fi FileInfo) error {
	if fi.Name == "Gsync" {
		fi.Id = g.remoteGsyncDir
		fi.Path = "Gsync"
	}

	files, err := g.store.ListFilesInDirectory(fi.Path, fi.Id)
	if err != nil {
		return err
	}

	for _, file := range files {
		if g.store.IsDir(file) {
			err = createDir(g.localGsyncDir + "/" + file.Name)
			if err != nil {
				return err
			}

			if err = g.Pull(file); err != nil {
				return err
			}

			continue
		}

		file.Data, err = g.store.GetFile(file)
		if err != nil {
			return err
		}

		err = os.WriteFile(g.localGsyncDir+"/"+file.Name, file.Data, 0700)
		if err != nil {
			return err
		}
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
