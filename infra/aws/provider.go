package aws

import (
	syncer2 "deni1688/gsync/syncer"
)

type awsS3Provider struct {
}

func NewSyncProvider() syncer2.SyncProvider {
	return &awsS3Provider{}
}

func (a awsS3Provider) IsDir(syncFile syncer2.SyncFile) bool {
	//TODO implement me
	panic("implement me")
}

func (a awsS3Provider) GetFile(syncFile syncer2.SyncFile) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (a awsS3Provider) CreateDir(syncFile syncer2.SyncFile) (syncer2.SyncFile, error) {
	//TODO implement me
	panic("implement me")
}

func (a awsS3Provider) CreateFile(syncFile syncer2.SyncFile) (syncer2.SyncFile, error) {
	//TODO implement me
	panic("implement me")
}

func (a awsS3Provider) UpdateFile(syncFile syncer2.SyncFile) error {
	//TODO implement me
	panic("implement me")
}

func (a awsS3Provider) ListFiles(dir syncer2.SyncFile) ([]syncer2.SyncFile, error) {
	//TODO implement me
	panic("implement me")
}
