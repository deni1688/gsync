package main

import (
	"deni1688/gsync/domain/syncer"
	"deni1688/gsync/infrastructure/cobraCli"
	"deni1688/gsync/infrastructure/googleDrive"
	"log"
	"os"
)

func main() {
	creds := os.Getenv("GOOGLE_OAUTH_CREDENTIALS_FILE")
	localDir := os.Getenv("LOCAL_GSYNC_DIR")

	gd := googleDrive.NewDrive(creds)
	ss := syncer.NewService(localDir, gd)
	rt := cobraCli.NewRuntime(ss)

	if err := rt.Execute(); err != nil {
		log.Fatalf("Error starting the CLI runtime: %v", err)
	}
}
