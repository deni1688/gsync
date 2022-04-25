package main

import (
	"deni1688/gsync/domain"
	"deni1688/gsync/infra/cli"
	"deni1688/gsync/infra/googledrive"
	"log"
)

func main() {
	googleDriverStore := googledrive.New()
	gsyncService := domain.NewGsyncService(googleDriverStore)
	runtime := cli.New(gsyncService)

	if err := runtime.Execute(); err != nil {
		log.Fatalf("Error starting the CLI runtime: %v", err)
	}
}
