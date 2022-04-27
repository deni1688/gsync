package domain

import (
	"log"
	"os"
)

type gsyncService struct {
	remoteGsyncDir string
	localGsyncDir  string
	store          SynchronizableStore
}

func NewGsyncService(localGsyncDir string, store SynchronizableStore) GsyncService {
	if localGsyncDir == "" {
		localGsyncDir = os.Getenv("HOME") + "/Gsync"
	}

	err := createDir(localGsyncDir)
	if err != nil {
		log.Fatalf("Error creating local gsync directory: %v", err)
	}

	dirFileInfo := FileInfo{
		Name: "Gsync",
		Path: localGsyncDir,
	}

	dirFileInfo, err = store.CreateDir(dirFileInfo)
	if err != nil {
		log.Fatalf("Error creating remote Gsync directory: %v", err)
	}

	return &gsyncService{dirFileInfo.Id, localGsyncDir, store}
}

func (g gsyncService) Sync(info FileInfo) error {
	var err error

	err = g.Push(info)
	err = g.Pull(info)

	return err
}
