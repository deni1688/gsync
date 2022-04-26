package domain

type SyncPriority int

const (
	SyncPriorityRemote SyncPriority = iota
	SyncPriorityLocal
)

type SyncOption struct {
	FilePath string
	Priority SyncPriority
}

type FileInfo struct {
	Name     string
	Size     int64
	MimeType string
	ParentId string
}
