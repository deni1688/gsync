package domain

type SynchronizableStoreContract interface {
	GetFile(info FileInfo) ([]byte, error)
	CreateFile(info FileInfo, data []byte) (FileInfo, error)
	CreateDirectory(name string) (FileInfo, error)
	UpdateFile(info FileInfo, data []byte) error
	ListFiles(root, directory string) ([]FileInfo, error)
	IsDir(info FileInfo) bool
}

type GsyncServiceContract interface {
	Pull(info FileInfo) error
	Push(info FileInfo) error
	Sync(info FileInfo) error
}
