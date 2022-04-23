package main

import (
	"deni1688/gsync/drive"
	"deni1688/gsync/store"
	"os"
)

func main() {
	home := os.Getenv("HOME")

	localPath := home + "/Gsync"
	localConfigPath := home + "/.gsync"

	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		os.Mkdir(localPath, 0700)
	}

	if _, err := os.Stat(localConfigPath); os.IsNotExist(err) {
		os.Mkdir(localConfigPath, 0700)
	}

	service := drive.New()
	store := store.New(service, localPath)

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
