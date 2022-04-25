package main

import (
	"deni1688/gsync/infra/drive"
	"deni1688/gsync/infra/store"
	"log"
	"os"
)

func main() {
	home := os.Getenv("HOME")
	localPath := home + "/Gsync"
	localConfigPath := home + "/.gsync"

	checkExists(localPath)
	checkExists(localConfigPath)

	service := drive.New()
	gsyncStore := store.New(service, localPath)

	handleCommands(gsyncStore, localPath)
}

func handleCommands(store *store.GsyncStore, localPath string) {
	switch getCommand() {
	case "pull":
		store.Pull(localPath, store.RemoteRoot.Id)
	case "push":
		store.Push(localPath, store.RemoteRoot.Id)
	default: // do nothing
	}
}

func getCommand() string {
	args := os.Args

	if len(args) < 2 {
		return ""
	}

	return args[1]
}

func checkExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, 0700); err != nil {
			log.Fatalf("Unable to create directory: %v", err)
		}
	}
}
