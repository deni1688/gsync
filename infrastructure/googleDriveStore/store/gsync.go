package store

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

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
		log.Fatalf("Could not read dir %s: %v\n", root, err)
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
			log.Fatalf("Could not open %s: %v\n", path, err)
		}
		defer media.Close()

		file := g.GetOrCreateRemote(f.Name(), false, parentId)
		file, err = g.Service.Files.Update(file.Id, &drive.File{Name: file.Name, MimeType: file.MimeType}).Media(media).Do()
		if err != nil {
			log.Fatalf("Could not update file %s: %v\n", file.Name, err)
		}
	}
}

func (g *GsyncStore) Pull(root, parentId string) {
	q := fmt.Sprintf("'%s' in parents and trashed = false", parentId)
	list, err := g.Service.Files.List().Q(q).Do()
	if err != nil {
		log.Fatalf("Could not fetch file list: %v\n", err)
	}

	for _, f := range list.Files {
		fileFullPath := fmt.Sprintf("/%s/%s", root, f.Name)
		if f.MimeType == "application/vnd.google-apps.folder" {
			if _, err := os.Stat(fileFullPath); os.IsNotExist(err) {
				os.Mkdir(fileFullPath, 0700)
			}

			g.Pull(fileFullPath, f.Id)
			continue
		}

		if strings.Contains(f.MimeType, "google-apps") {
			continue
		}

		resp, err := g.Service.Files.Get(f.Id).Download()
		if err != nil {
			log.Fatalf("Could not fetch file %s in %s: %v\n", f.Name, root, err)
		}
		defer resp.Body.Close()

		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Could read file %s in %s: %v\n", f.Name, root, err)
		}

		err = os.WriteFile(fileFullPath, bytes, 0700)
		if err != nil {
			log.Fatalf("Could write file %s in %s: %v\n", f.Name, root, err)
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

	q := fmt.Sprintf("name = '%s' and trashed = false", name)

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
