package main

import (
	"deni1688/gsync/domain"
	"deni1688/gsync/infra/cli"
	"deni1688/gsync/infra/googledrive"
	"log"
)

func main() {
	googleDriverStorage := googledrive.NewGoogleDriveStorage()
	gsyncService := domain.NewGsyncService(googleDriverStorage)
	cliRuntime := cli.NewCliRuntime(gsyncService)

	if err := cliRuntime.Execute(); err != nil {
		log.Fatalf("Error starting the CLI runtime: %v", err)
	}
}
