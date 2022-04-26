package main

import (
	"deni1688/gsync/domain"
	"deni1688/gsync/infrastructure/cli"
	"deni1688/gsync/infrastructure/googleDriveStore"
	"log"
	"os"
)

func main() {
	credentialsPath := os.Getenv("GOOGLE_OAUTH_CREDENTIALS")
	localGsyncDir := os.Getenv("LOCAL_GSYNC_DIR")

	googleDriverStore := googleDriveStore.NewSecurityKey(credentialsPath)
	gsyncService := domain.NewGsyncService(localGsyncDir, googleDriverStore)
	runtime := cli.New(gsyncService)

	if err := runtime.Execute(); err != nil {
		log.Fatalf("Error starting the CLI runtime: %v", err)
	}
}
