package main

import (
	"deni1688/gsync/domain"
	"deni1688/gsync/infra/aws"
	"deni1688/gsync/infra/cobra"
	"deni1688/gsync/infra/google"
	"log"
	"os"
)

var (
	creds          = os.Getenv("GOOGLE_OAUTH_CREDENTIALS_FILE")
	localDir       = os.Getenv("LOCAL_GSYNC_DIR")
	providerConfig = os.Getenv("SYNC_PROVIDER")
)

func main() {
	sp := selectSyncProvider(providerConfig)
	gs := domain.NewService(localDir, sp)

	if err := cobra.NewCLI(gs).Execute(); err != nil {
		log.Fatalf("Error starting the CLI runtime: %v", err)
	}
}

func selectSyncProvider(providerConfig string) domain.SyncProvider {
	if providerConfig == "" {
		log.Fatalf("SYNC_PROVIDER env variable required! You can specify  google or aws")
	}

	var sp domain.SyncProvider
	if providerConfig == "google" {
		sp = google.NewSyncProvider(creds)
	} else {
		sp = aws.NewSyncProvider()
	}
	return sp
}
