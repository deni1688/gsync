package main

import (
	"deni1688/gsync/drive"
	"deni1688/gsync/store"
	"os"
)

func main() {
	home := os.Getenv("HOME")
	localPath := home + "/Gsync"

	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		os.Mkdir(localPath, 0700)
	}

	service := drive.New()
	store := store.New(service, localPath)
	store.Push(store.LocalRoot, store.RemoteRoot.Id)
}
