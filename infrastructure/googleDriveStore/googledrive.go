package googleDriveStore

import (
	"bytes"
	"context"
	"deni1688/gsync/domain"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"log"
	"os"
)

type store struct {
	service *drive.Service
}

func (s store) GetFile(info domain.FileInfo) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (s store) CreateFile(info domain.FileInfo, data []byte) error {
	file := &drive.File{Name: info.Name, MimeType: info.MimeType}

	if info.ParentId != "" {
		file.Parents = []string{info.ParentId}
	}

	q := fmt.Sprintf("name = '%s' and trashed = false", info.Name)

	list, err := s.service.Files.List().Fields("files(id, name, mimeType)").Q(q).Do()
	if err != nil {
		log.Printf("Could not fetch list: %v", err)
		return err
	}

	if len(list.Files) == 0 {
		fmt.Printf("Could not find %s %s...creating\n", info.MimeType, info.Name)
		file, err = s.service.Files.Create(file).Do()

		if err != nil {
			return err
		}
	}

	fmt.Printf("Found %s %s...upadting\n", info.MimeType, info.Name)
	file = list.Files[0]
	file, err = s.service.Files.Update(file.Id, &drive.File{Name: file.Name, MimeType: file.MimeType}).Media(bytes.NewReader(data)).Do()
	if err != nil {
		log.Printf("Could update file: %v", err)
		return err
	}

	return nil
}

func (s store) UpdateFile(info domain.FileInfo, data []byte) error {
	//TODO implement me
	panic("implement me")
}

func (s store) DeleteFile(info domain.FileInfo) error {
	//TODO implement me
	panic("implement me")
}

func (s store) ListFiles(path string) ([]domain.FileInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (s store) IsDir(info domain.FileInfo) bool {
	return info.MimeType == "application/vnd.google-apps.folder"
}

func New(credentialsPath string) domain.SynchronizableStoreContract {
	b, err := os.ReadFile(credentialsPath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := getClient(config)

	service, err := drive.NewService(context.TODO(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	return &store{service}
}
