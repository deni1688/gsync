package main

import (
	"deni1688/gsync/domain/syncer"
	"deni1688/gsync/infrastructure/aws"
	"deni1688/gsync/infrastructure/cobra"
	"deni1688/gsync/infrastructure/google"
	"log"
	"os"
)

var (
	creds          = os.Getenv("GOOGLE_OAUTH_CREDENTIALS_FILE")
	localDir       = os.Getenv("LOCAL_GSYNC_DIR")
	providerConfig = os.Getenv("SYNC_PROVIDER")
)

func main() {
	sp := getSyncProvider(providerConfig)

	gs := syncer.NewService(localDir, sp)
	c := cobra.NewCLI(gs)

	if err := c.Execute(); err != nil {
		log.Fatalf("Error starting the CLI runtime: %v", err)
	}
}

func getSyncProvider(providerConfig string) syncer.SyncProvider {
	if providerConfig == "" {
		log.Fatalf("SYNC_PROVIDER env variable required! You can specify  google or aws")
	}

	var sp syncer.SyncProvider
	if providerConfig == "google" {
		sp = google.NewSyncProvider(creds)
	} else {
		sp = aws.NewSyncProvider()
	}
	return sp
}
