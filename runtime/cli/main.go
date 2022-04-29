package main

import (
	"deni1688/gsync/domain/syncer"
	"deni1688/gsync/infrastructure/cobracli"
	"deni1688/gsync/infrastructure/googledrive"
	"log"
	"os"
)

func main() {
	creds := os.Getenv("GOOGLE_OAUTH_CREDENTIALS_FILE")
	localDir := os.Getenv("LOCAL_GSYNC_DIR")

	gd := googledrive.NewDrive(creds)
	ss := syncer.NewService(localDir, gd)
	rt := cobracli.NewRuntime(ss)

	if err := rt.Execute(); err != nil {
		log.Fatalf("Error starting the CLI runtime: %v", err)
	}
}
