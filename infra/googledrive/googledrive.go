package googledrive

import "deni1688/gsync/domain"

type GoogleDriveStorage struct {
}

func (g GoogleDriveStorage) Authorize() error {
	//TODO implement me
	panic("implement me")
}

func NewGoogleDriveStorage() domain.SynchronizableStorageContract {
	return &GoogleDriveStorage{}
}

func (g GoogleDriveStorage) Pull(option ...domain.SyncOption) error {
	//TODO implement me
	panic("implement me")
}

func (g GoogleDriveStorage) Push(option ...domain.SyncOption) error {
	//TODO implement me
	panic("implement me")
}

func (g GoogleDriveStorage) Sync(option ...domain.SyncOption) error {
	//TODO implement me
	panic("implement me")
}
