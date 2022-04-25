package googledrive

import (
	"deni1688/gsync/domain"
)

type store struct {
}

func (s store) GetAuthorizationToken(credentialsPath string) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (s store) GetFile(info domain.FileInfo) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (s store) CreateFile(info, FileInfo, data []byte) error {
	//TODO implement me
	panic("implement me")
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

func New() domain.SynchronizableStoreContract {
	return &store{}
}
