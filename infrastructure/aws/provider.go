package aws

import (
	syncer2 "deni1688/gsync/domain"
)

type awsS3Provider struct {
}

func NewSyncProvider() syncer2.SyncProvider {
	return &awsS3Provider{}
}

func (a awsS3Provider) IsDir(syncTarget syncer2.SyncTarget) bool {
	//TODO implement me
	panic("implement me")
}

func (a awsS3Provider) GetFile(syncTarget syncer2.SyncTarget) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (a awsS3Provider) CreateDir(syncTarget syncer2.SyncTarget) (syncer2.SyncTarget, error) {
	//TODO implement me
	panic("implement me")
}

func (a awsS3Provider) CreateFile(syncTarget syncer2.SyncTarget) (syncer2.SyncTarget, error) {
	//TODO implement me
	panic("implement me")
}

func (a awsS3Provider) UpdateFile(syncTarget syncer2.SyncTarget) error {
	//TODO implement me
	panic("implement me")
}

func (a awsS3Provider) ListFiles(dir syncer2.SyncTarget) ([]syncer2.SyncTarget, error) {
	//TODO implement me
	panic("implement me")
}
