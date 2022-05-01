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

type gsyncService struct {
	remoteGsyncDir string
	localGsyncDir  string
	syncProvider   SyncProvider
}

func NewService(localGsyncDir string, syncProvider SyncProvider) GsyncService {
	if localGsyncDir == "" {
		localGsyncDir = os.Getenv("HOME") + "/Gsync"
	}

	err := createDirs(localGsyncDir, getPath(os.Getenv("HOME"), ".gsync"))
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

	return &gsyncService{dir.Id, localGsyncDir, syncProvider}
}

func (gs gsyncService) Sync(dir SyncFile) error {
	if err := gs.Push(dir); err != nil {
		return err
	}

	if err := gs.Pull(dir); err != nil {
		return err
	}

	return nil
}
