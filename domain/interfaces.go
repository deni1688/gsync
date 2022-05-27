package domain

type SyncProvider interface {
	IsDir(syncTarget SyncTarget) bool
	GetFile(syncTarget SyncTarget) ([]byte, error)
	CreateDir(syncTarget SyncTarget) (SyncTarget, error)
	CreateFile(syncTarget SyncTarget) (SyncTarget, error)
	UpdateFile(syncTarget SyncTarget) error
	ListFiles(dir SyncTarget) ([]SyncTarget, error)
}

type GsyncService interface {
	Pull(dir SyncTarget) error
	Push(dir SyncTarget) error
	Sync(dir SyncTarget) error
}
