package syncer

import (
	"log"
	"os"
	"sync"
	"time"
)

func (g syncService) Pull(sf SyncFile) error {
	if sf.Name == "Gsync" {
		sf.Id = g.remoteGsyncDir
		sf.Path = g.localGsyncDir
	}

	files, err := g.drive.ListFiles(sf)
	if err != nil {
		return err
	}

	if err = g.removeFilesFromLocal(sf, files); err != nil {
		return err
	}

	return g.downloadFiles(sf, files)
}

func (g syncService) removeFilesFromLocal(sf SyncFile, files []SyncFile) error {
	list, err := os.ReadDir(sf.Path)
	if err != nil {
		return err
	}

	for _, file := range list {
		name := file.Name()
		fullPath := GetFullPath(sf.Path, name)

		if SyncFileListContains(files, name) {
			continue
		}

		log.Printf("Removing %s", fullPath)

		if file.IsDir() {
			if err = os.RemoveAll(fullPath); err != nil {
				return err
			}
		} else {
			if err = os.Remove(fullPath); err != nil {
				return err
			}
		}
	}

	return err
}

func (g syncService) downloadFiles(sf SyncFile, files []SyncFile) error {
	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	for _, file := range files {
		fullPath := GetFullPath(sf.Path, file.Name)

		wg.Add(1)
		go func(wg *sync.WaitGroup, file SyncFile) {
			defer wg.Done()

			if g.drive.IsDir(file) {
				log.Printf("Pulling dir %s", fullPath)

				if err := CreateDir(fullPath); err != nil {
					errCh <- err
				}

				file.Path = fullPath

				if err := g.Pull(file); err != nil {
					errCh <- err
				}
			} else {
				log.Printf("Pulling file %s", fullPath)
				data, err := g.drive.GetFile(file)
				if err != nil {
					errCh <- err
				}

				if err = os.WriteFile(fullPath, data, 0700); err != nil {
					errCh <- err
				}
			}
		}(&wg, file)

		time.Sleep(time.Millisecond * 30) // throttle api calls
	}
	wg.Wait()

	if len(errCh) < 1 {
		errCh <- nil
	}

	return <-errCh
}
