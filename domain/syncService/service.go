package syncService

import (
	"deni1688/gsync/domain"
	"log"
	"os"
)

type gsyncService struct {
	remoteGsyncDir string
	localGsyncDir  string
	drive          domain.SynchronizableDrive
}

func New(localGsyncDir string, drive domain.SynchronizableDrive) domain.GsyncService {
	if localGsyncDir == "" {
		localGsyncDir = os.Getenv("HOME") + "/Gsync"
	}

	err := domain.CreateDir(localGsyncDir)
	if err != nil {
		log.Fatalf("Error creating local gsync directory: %v", err)
	}

	dirSyncFile := domain.SyncFile{
		Name: "Gsync",
		Path: localGsyncDir,
	}

	dirSyncFile, err = drive.CreateDir(dirSyncFile)
	if err != nil {
		log.Fatalf("Error creating remote Gsync directory: %v", err)
	}

	return &gsyncService{dirSyncFile.Id, localGsyncDir, drive}
}

func (g gsyncService) Sync(syncFile domain.SyncFile) error {
	var err error

	err = g.Push(syncFile)
	err = g.Pull(syncFile)

	return err
}
