package syncer

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

func (g syncService) PullFiles(dir SyncFile) error {
	if dir.Name == "Gsync" {
		dir.Id = g.remoteGsyncDir
		dir.Path = g.localGsyncDir
	}

	if dir.Id == "" {
		return fmt.Errorf("dir id is required")
	}

	files, err := g.syncProvider.ListFiles(dir)
	if err != nil {
		return err
	}

	if err = g.cleanLocalFiles(dir, files); err != nil {
		return err
	}

	return g.downloadFiles(dir, files)
}

func (g syncService) cleanLocalFiles(dir SyncFile, files []SyncFile) error {
	list, err := os.ReadDir(dir.Path)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	for _, file := range list {
		name := file.Name()
		fullPath := GetPathFrom(dir.Path, name)

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

func (g syncService) downloadFiles(dir SyncFile, files []SyncFile) error {
	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	for _, file := range files {
		fullPath := GetPathFrom(dir.Path, file.Name)

		wg.Add(1)
		go func(wg *sync.WaitGroup, file SyncFile) {
			defer wg.Done()

			if g.syncProvider.IsDir(file) {
				log.Printf("Pulling dir %s", fullPath)

				if err := CreateDir(fullPath); err != nil {
					errCh <- fmt.Errorf("could not create dir %s: %v\n", fullPath, err)
				}

				file.Path = fullPath

				if err := g.PullFiles(file); err != nil {
					errCh <- fmt.Errorf("could not pull dir %s: %v\n", fullPath, err)
				}
			} else {
				log.Printf("Pulling file %s", fullPath)

				data, err := g.syncProvider.GetFile(file)
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
