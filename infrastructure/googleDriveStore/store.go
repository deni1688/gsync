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
	"net/http"
	"os"
	"strings"
)

var mimeTypeMap = map[string]string{
	"application/vnd.google-apps.file":         "text/plain",
	"application/vnd.google-apps.drawing":      "image/png",
	"application/vnd.google-apps.photo":        "image/jpeg",
	"application/vnd.google-apps.video":        "video/mp4",
	"application/vnd.google-apps.audio":        "audio/mpeg",
	"application/vnd.google-apps.spreadsheet":  "application/x-vnd.oasis.opendocument.spreadsheet",
	"application/vnd.google-apps.document":     "application/vnd.oasis.opendocument.text",
	"application/vnd.google-apps.presentation": "application/vnd.oasis.opendocument.presentation",
}

type store struct {
	service *drive.Service
}

func (s store) CreateDirectory(name string) (domain.FileInfo, error) {
	q := fmt.Sprintf("name = '%s' and trashed = false", name)
	list, err := s.service.Files.List().Fields("files(id, name, mimeType)").Q(q).Do()
	if err != nil {
		return domain.FileInfo{}, err
	}

	if len(list.Files) > 0 {
		f := list.Files[0]
		return domain.FileInfo{
			Id:       f.Id,
			Name:     f.Name,
			MimeType: f.MimeType,
			ParentId: getParentId(f),
			Path:     f.Name,
		}, nil
	}

	file := &drive.File{Name: name, MimeType: "application/vnd.google-apps.folder"}
	file, err = s.service.Files.Create(file).Do()
	if err != nil {
		return domain.FileInfo{}, err
	}

	return domain.FileInfo{
		Id:       file.Id,
		Name:     file.Name,
		MimeType: file.MimeType,
		ParentId: getParentId(file),
		Path:     file.Name,
	}, nil
}

func getParentId(f *drive.File) string {
	if len(f.Parents) > 0 {
		return f.Parents[0]
	}
	return ""
}

func (s store) GetFile(info domain.FileInfo) ([]byte, error) {
	var err error
	resp := new(http.Response)

	if strings.Contains(info.MimeType, "google-apps") {
		resp, err = s.service.Files.Export(info.Id, getExportMimeType(info.MimeType)).Download()
	} else {
		resp, err = s.service.Files.Get(info.Id).Download()
	}

	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	return io.ReadAll(resp.Body)
}

func getExportMimeType(driveMimeType string) string {
	if mimeType, ok := mimeTypeMap[driveMimeType]; ok {
		return mimeType
	}
	return "application/octet-stream"
}

func (s store) CreateFile(info domain.FileInfo, data []byte) (domain.FileInfo, error) {
	file := &drive.File{Name: info.Name, MimeType: info.MimeType}

	if info.ParentId != "" {
		file.Parents = []string{info.ParentId}
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

func (s store) ListFiles(root, parentId string) ([]domain.FileInfo, error) {
	q := fmt.Sprintf("'%s' in parents and trashed = false", parentId)
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
			ParentId: getParentId(f),
			Path:     root + "/" + f.Name,
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
