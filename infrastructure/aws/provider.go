package aws

import (
	"deni1688/gsync/domain/syncer"
)

type awsS3Provider struct {
}

func NewSyncProvider() syncer.SyncProvider {
	return &awsS3Provider{}
}

func (a awsS3Provider) IsDir(syncFile syncer.SyncFile) bool {
	//TODO implement me
	panic("implement me")
}

func (a awsS3Provider) GetFile(syncFile syncer.SyncFile) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (a awsS3Provider) CreateDir(syncFile syncer.SyncFile) (syncer.SyncFile, error) {
	//TODO implement me
	panic("implement me")
}

func (a awsS3Provider) CreateFile(syncFile syncer.SyncFile) (syncer.SyncFile, error) {
	//TODO implement me
	panic("implement me")
}

func (a awsS3Provider) UpdateFile(syncFile syncer.SyncFile) error {
	//TODO implement me
	panic("implement me")
}

func (a awsS3Provider) ListFiles(dir syncer.SyncFile) ([]syncer.SyncFile, error) {
	//TODO implement me
	panic("implement me")
}
