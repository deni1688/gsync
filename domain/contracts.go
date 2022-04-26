package domain

type SynchronizableStoreContract interface {
	GetFile(info FileInfo) ([]byte, error)
	CreateFile(info FileInfo, data []byte) error
	UpdateFile(info FileInfo, data []byte) error
	DeleteFile(info FileInfo) error
	ListFiles(path string) ([]FileInfo, error)
	IsDir(info FileInfo) bool
}

type GsyncServiceContract interface {
	Pull(path string) error
	Push(option ...SyncOption) error
	Sync(option ...SyncOption) error
}
