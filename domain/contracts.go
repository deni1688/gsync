package domain

type SynchronizableStoreContract interface {
	GetFile(info FileInfo) ([]byte, error)
	CreateFile(info FileInfo) (FileInfo, error)
	CreateDir(info FileInfo) (FileInfo, error)
	UpdateFile(info FileInfo) error
	ListFiles(parentFileInfo FileInfo) ([]FileInfo, error)
	IsDir(info FileInfo) bool
}

type GsyncServiceContract interface {
	Pull(info FileInfo) error
	Push(info FileInfo) error
	Sync(info FileInfo) error
}
