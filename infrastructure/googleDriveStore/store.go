package googleDriveStore

import (
	"bytes"
	"context"
	"deni1688/gsync/domain"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"io"
	"log"
	"os"
)

type store struct {
	service *drive.Service
}

func (s store) CreateDirectory(name string) (domain.FileInfo, error) {
	file := &drive.File{Name: name, MimeType: "application/vnd.google-apps.folder"}

	file, err := s.service.Files.Create(file).Do()
	if err != nil {
		return domain.FileInfo{}, err
	}

	return domain.FileInfo{
		Id:       file.Id,
		Name:     file.Name,
		MimeType: file.MimeType,
	}, nil
}

func (s store) GetFile(info domain.FileInfo) ([]byte, error) {
	resp, err := s.service.Files.Get(info.Id).Download()
	if err != nil {
		log.Fatalf("Could not fetch file %s: %v", info.Name, err)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (s store) CreateFile(info domain.FileInfo, data []byte) (domain.FileInfo, error) {
	file := &drive.File{Name: info.Name, MimeType: info.MimeType}

	if info.Parent != "" {
		file.Parents = []string{info.Parent}
	}

	file, err := s.service.Files.Create(file).Do()
	if err != nil {
		log.Fatalf("Could not create file %s: %v", info.Name, err)
	}

	info.Id = file.Id

	return info, s.UpdateFile(info, data)
}

func (s store) UpdateFile(info domain.FileInfo, data []byte) error {
	file := &drive.File{Id: info.Id, Name: info.Name, MimeType: info.MimeType}

	file, err := s.service.Files.Update(file.Id, &drive.File{Name: file.Name, MimeType: file.MimeType}).Media(bytes.NewReader(data)).Do()
	if err != nil {
		log.Printf("Could update file: %v", err)
		return err
	}

	return nil
}

func (s store) ListFiles(path string) ([]domain.FileInfo, error) {
	q := fmt.Sprintf("name = '%s' and trashed = false", path)

	list, err := s.service.Files.List().Fields("files(id, name, mimeType)").Q(q).Do()
	if err != nil {
		return nil, err
	}

	var files []domain.FileInfo
	for _, f := range list.Files {
		files = append(files, domain.FileInfo{
			Id:       f.Id,
			Name:     f.Name,
			MimeType: f.MimeType,
		})
	}

	return files, nil
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
