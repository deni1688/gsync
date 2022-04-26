package domain

type SynchronizableStoreContract interface {
	GetFile(info FileInfo) ([]byte, error)
	CreateFile(info FileInfo, data []byte) (FileInfo, error)
	CreateDirectory(name string) (FileInfo, error)
	UpdateFile(info FileInfo, data []byte) error
	ListFiles(path string) ([]FileInfo, error)
	IsDir(info FileInfo) bool
}

type GsyncServiceContract interface {
	Pull(path string) error
	Push(option ...SyncOption) error
	Sync(option ...SyncOption) error
}
