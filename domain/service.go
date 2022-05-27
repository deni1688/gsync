package domain

import (
	"log"
	"os"
)

type SyncTarget struct {
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

	dir := SyncTarget{
		Name: "Gsync",
		Path: localGsyncDir,
	}

	dir, err = syncProvider.CreateDir(dir)
	if err != nil {
		log.Fatalf("Error creating remote Gsync directory: %v", err)
	}

	return &gsyncService{dir.Id, localGsyncDir, syncProvider}
}
