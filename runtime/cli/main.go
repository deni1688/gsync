package main

import (
	"deni1688/gsync/domain/synchronizer"
	"deni1688/gsync/infrastructure/cobraCliRuntime"
	"deni1688/gsync/infrastructure/googleDrive"
	"log"
	"os"
)

func main() {
	creds := os.Getenv("GOOGLE_OAUTH_CREDENTIALS_FILE")
	localDir := os.Getenv("LOCAL_GSYNC_DIR")

	gd := googleDrive.New(creds)
	ss := synchronizer.New(localDir, gd)
	rt := cobraCliRuntime.New(ss)

	if err := rt.Execute(); err != nil {
		log.Fatalf("Error starting the CLI runtime: %v", err)
	}
}
