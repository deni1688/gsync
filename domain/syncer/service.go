package syncer

import (
	"log"
	"os"
)

type syncService struct {
	remoteGsyncDir string
	localGsyncDir  string
	drive          SynchronizableDrive
}

func NewService(localGsyncDir string, drive SynchronizableDrive) GsyncService {
	if localGsyncDir == "" {
		localGsyncDir = os.Getenv("HOME") + "/Gsync"
	}

	err := CreateDir(localGsyncDir)
	if err != nil {
		log.Fatalf("Error creating local gsync directory: %v", err)
	}

	dirSyncFile := SyncFile{
		Name: "Gsync",
		Path: localGsyncDir,
	}

	dirSyncFile, err = drive.CreateDir(dirSyncFile)
	if err != nil {
		log.Fatalf("Error creating remote Gsync directory: %v", err)
	}

	return &syncService{dirSyncFile.Id, localGsyncDir, drive}
}

func (g syncService) Sync(syncFile SyncFile) error {
	var err error

	err = g.Push(syncFile)
	err = g.Pull(syncFile)

	return err
}
