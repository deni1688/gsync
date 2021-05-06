package store

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"google.golang.org/api/drive/v3"
)

type GsyncStore struct {
	Service    *drive.Service
	RemoteRoot *drive.File
	LocalRoot  string
}

func (g *GsyncStore) Push(root, parentId string) {
	files, err := ioutil.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(f.ModTime())
		path := root + "/" + f.Name()
		fmt.Println(f.Name())

		if f.IsDir() {
			dir := g.GetOrCreateRemote(f.Name(), true, parentId)
			g.Push(path, dir.Id)
			continue
		}

		media, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		defer media.Close()

		file := g.GetOrCreateRemote(f.Name(), false, parentId)
		file, err = g.Service.Files.Update(file.Id, &drive.File{Name: file.Name, MimeType: file.MimeType}).Media(media).Do()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (g *GsyncStore) GetOrCreateRemote(name string, isDir bool, parentId string) *drive.File {
	file := &drive.File{Name: name}

	t := "file"
	if isDir {
		t = "folder"
		file.MimeType = "application/vnd.google-apps.folder"
	}

	if parentId != "" {
		file.Parents = []string{parentId}
	}

	q := fmt.Sprintf("name = '%s'", name)

	list, err := g.Service.Files.List().Fields("files(id, name, mimeType)").Q(q).Do()
	if err != nil {
		log.Fatalf("Could not fetch list: %v", err)
	}
	if len(list.Files) == 0 {
		fmt.Printf("Could not find %s %s...creating\n", t, name)
		file, _ = g.Service.Files.Create(file).Do()
	} else {
		fmt.Printf("Found %s %s...upadting\n", t, name)
		file = list.Files[0]
	}
	return file
}

func New(service *drive.Service, localPath string) *GsyncStore {
	store := new(GsyncStore)
	store.Service = service
	store.RemoteRoot = new(drive.File)
	store.LocalRoot = localPath

	store.RemoteRoot = store.GetOrCreateRemote("Gsync", true, "")

	return store
}
