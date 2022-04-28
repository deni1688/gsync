package syncer

type SynchronizableDrive interface {
	GetFile(syncFile SyncFile) ([]byte, error)
	CreateFile(syncFile SyncFile) (SyncFile, error)
	CreateDir(syncFile SyncFile) (SyncFile, error)
	UpdateFile(syncFile SyncFile) error
	ListFiles(parentSyncFile SyncFile) ([]SyncFile, error)
	IsDir(syncFile SyncFile) bool
}

type GsyncService interface {
	Pull(syncFile SyncFile) error
	Push(syncFile SyncFile) error
	Sync(syncFile SyncFile) error
}
