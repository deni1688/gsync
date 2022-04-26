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
	Id       string
	Name     string
	MimeType string
	Parent   string
}
