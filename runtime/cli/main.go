package main

import (
	"deni1688/gsync/domain"
	"deni1688/gsync/infra/cli"
	"deni1688/gsync/infra/googleDriveStore"
	"log"
	"os"
)

func main() {
	googleDriverStore := googleDriveStore.New()
	gsyncService := domain.NewGsyncService(os.Getenv("LOCAL_GSYNC_DIR"), googleDriverStore)
	runtime := cli.New(gsyncService)

	if err := runtime.Execute(); err != nil {
		log.Fatalf("Error starting the CLI runtime: %v", err)
	}
}
