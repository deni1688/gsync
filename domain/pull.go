package domain

import (
	"fmt"
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

	if err = g.cleanLocalFiles(sf, files); err != nil {
		return err
	}

	return g.downloadFiles(sf, files)
}

func (g syncService) cleanLocalFiles(sf SyncFile, files []SyncFile) error {
	list, err := os.ReadDir(sf.Path)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	for _, file := range list {
		name := file.Name()
		fullPath := GetPathFrom(sf.Path, name)

		wg.Add(1)
		go func(wg *sync.WaitGroup, file os.DirEntry) {
			if FileListContains(files, name) {
				wg.Done()
				return
			}

			log.Printf("Removing %s", fullPath)

			if file.IsDir() {
				log.Printf("Removing dir %s", fullPath)
				if err = os.RemoveAll(fullPath); err != nil {
					errCh <- fmt.Errorf("could not remove dir %s: %v\n", fullPath, err)
				}
			} else {
				log.Printf("Removing file %s", fullPath)
				if err = os.Remove(fullPath); err != nil {
					errCh <- fmt.Errorf("could not remove file %s: %v\n", fullPath, err)
				}
			}
			wg.Done()
		}(&wg, file)
	}
	wg.Wait()

	if len(errCh) < 1 {
		errCh <- nil
	}

	return err
}

func (g syncService) downloadFiles(sf SyncFile, files []SyncFile) error {
	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	for _, file := range files {
		fullPath := GetPathFrom(sf.Path, file.Name)

		wg.Add(1)
		go func(wg *sync.WaitGroup, file SyncFile) {
			defer wg.Done()

			if g.drive.IsDir(file) {
				log.Printf("Pulling dir %s", fullPath)

				if err := CreateDir(fullPath); err != nil {
					errCh <- fmt.Errorf("could not create dir %s: %v\n", fullPath, err)
				}

				file.Path = fullPath

				if err := g.Pull(file); err != nil {
					errCh <- fmt.Errorf("could not pull dir %s: %v\n", fullPath, err)
				}
			} else {
				log.Printf("Pulling file %s", fullPath)

				data, err := g.drive.GetFile(file)
				if err != nil {
					errCh <- fmt.Errorf("could not get file %v: %v\n", file, err)
				}

				if err = os.WriteFile(fullPath, data, 0700); err != nil {
					errCh <- fmt.Errorf("could not write file %s: %v\n", fullPath, err)
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
