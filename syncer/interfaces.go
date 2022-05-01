package syncer

type SyncProvider interface {
	IsDir(syncFile SyncFile) bool
	GetFile(syncFile SyncFile) ([]byte, error)
	CreateDir(syncFile SyncFile) (SyncFile, error)
	CreateFile(syncFile SyncFile) (SyncFile, error)
	UpdateFile(syncFile SyncFile) error
	ListFiles(dir SyncFile) ([]SyncFile, error)
}

type GsyncService interface {
	Pull(dir SyncFile) error
	Push(dir SyncFile) error
	Sync(dir SyncFile) error
}
