package googleDrive

import (
	"bytes"
	"context"
	"deni1688/gsync/domain/syncer"
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

type googleDriveService struct {
	service *drive.Service
}

func (s googleDriveService) CreateDir(dirSyncFile syncer.SyncFile) (syncer.SyncFile, error) {
	q := fmt.Sprintf("name = '%s' and trashed = false", dirSyncFile.Name)

	list, err := s.service.Files.List().Fields("files(id, name, mimeType)").Q(q).Do()
	if err != nil {
		return syncer.SyncFile{}, err
	}

	if len(list.Files) > 0 {
		f := list.Files[0]
		return syncer.SyncFile{
			Id:       f.Id,
			Name:     f.Name,
			MimeType: f.MimeType,
			ParentId: getParentId(f),
			Path:     dirSyncFile.Path,
		}, nil
	}

	file := &drive.File{
		Name:     dirSyncFile.Name,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{dirSyncFile.ParentId},
	}

	file, err = s.service.Files.Create(file).Do()
	if err != nil {
		return syncer.SyncFile{}, err
	}

	return syncer.SyncFile{
		Id:       file.Id,
		Name:     file.Name,
		MimeType: file.MimeType,
		ParentId: getParentId(file),
		Path:     dirSyncFile.Path,
	}, nil
}

func getParentId(f *drive.File) string {
	if len(f.Parents) > 0 {
		return f.Parents[0]
	}
	return ""
}

func (s googleDriveService) GetFile(syncFile syncer.SyncFile) ([]byte, error) {
	var err error
	resp := new(http.Response)

	if strings.Contains(syncFile.MimeType, "google-apps") {
		resp, err = s.service.Files.Export(syncFile.Id, getExportMimeType(syncFile.MimeType)).Download()
	} else {
		resp, err = s.service.Files.Get(syncFile.Id).Download()
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

func (s googleDriveService) CreateFile(syncFile syncer.SyncFile) (syncer.SyncFile, error) {
	file := &drive.File{Name: syncFile.Name, MimeType: syncFile.MimeType}

	if syncFile.ParentId != "" {
		file.Parents = []string{syncFile.ParentId}
	}

	file, err := s.service.Files.Create(file).Media(bytes.NewReader(syncFile.Data)).Do()
	if err != nil {
		return syncer.SyncFile{}, err
	}

	syncFile.Id = file.Id

	return syncFile, nil
}

func (s googleDriveService) UpdateFile(syncFile syncer.SyncFile) error {
	file := &drive.File{Id: syncFile.Id, Name: syncFile.Name, MimeType: syncFile.MimeType}

	file, err := s.service.Files.Update(file.Id, &drive.File{Name: file.Name, MimeType: file.MimeType}).Media(bytes.NewReader(syncFile.Data)).Do()
	if err != nil {
		log.Printf("Could update file: %v", err)
		return err
	}

	return nil
}

func (s googleDriveService) ListFiles(parentSyncFile syncer.SyncFile) ([]syncer.SyncFile, error) {
	q := fmt.Sprintf("'%s' in parents and trashed = false", parentSyncFile.Id)

	list, err := s.service.Files.List().Fields("files(id, name, mimeType)").Q(q).Do()
	if err != nil {
		return nil, err
	}

	var files []syncer.SyncFile

	for _, f := range list.Files {
		files = append(files, syncer.SyncFile{
			Id:       f.Id,
			Name:     f.Name,
			MimeType: f.MimeType,
			ParentId: getParentId(f),
			Path:     parentSyncFile.Path + "/" + f.Name,
		})
	}

	return files, nil
}

func (s googleDriveService) IsDir(syncFile syncer.SyncFile) bool {
	return syncFile.MimeType == "application/vnd.google-apps.folder"
}

func NewDrive(credentialsPath string) syncer.SynchronizableDrive {
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

	return &googleDriveService{service}
}
