package domain

type SynchronizableStoreContract interface {
	GetFile(info FileInfo) ([]byte, error)
	CreateFile(info FileInfo, data []byte) (FileInfo, error)
	CreateDirectory(name string) (FileInfo, error)
	UpdateFile(info FileInfo, data []byte) error
	ListFilesInDirectory(root, directory string) ([]FileInfo, error)
	IsDir(info FileInfo) bool
}

type GsyncServiceContract interface {
	Pull(info FileInfo) error
	Push(option ...SyncOption) error
	Sync(option ...SyncOption) error
}
