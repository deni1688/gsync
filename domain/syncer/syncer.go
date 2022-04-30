package syncer

import (
	"log"
	"os"
)

type SyncFile struct {
	Id       string
	Name     string
	MimeType string
	ParentId string
	Path     string
	Data     []byte
}

type syncService struct {
	remoteGsyncDir string
	localGsyncDir  string
	syncProvider   SyncProvider
}

func NewService(localGsyncDir string, syncProvider SyncProvider) GsyncService {
	if localGsyncDir == "" {
		localGsyncDir = os.Getenv("HOME") + "/Gsync"
	}

	err := CreateDir(localGsyncDir)
	if err != nil {
		log.Fatalf("Error creating local gsync directory: %v", err)
	}

	dir := SyncFile{
		Name: "Gsync",
		Path: localGsyncDir,
	}

	dir, err = syncProvider.CreateDir(dir)
	if err != nil {
		log.Fatalf("Error creating remote Gsync directory: %v", err)
	}

	return &syncService{dir.Id, localGsyncDir, syncProvider}
}

func (g syncService) SyncFiles(dir SyncFile) error {
	if err := g.PushFiles(dir); err != nil {
		return err
	}

	if err := g.PullFiles(dir); err != nil {
		return err
	}

	return nil
}
