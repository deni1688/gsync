package domain

type SyncPriority int

const (
	SyncPriorityRemote SyncPriority = iota
	SyncPriorityLocal
)

type SyncOption struct {
	Filename string       `json:"filename"`
	Priority SyncPriority `json:"priority"`
}

type SynchronizableStorageContract interface {
	Pull(option ...SyncOption) error
	Push(option ...SyncOption) error
	Sync(option ...SyncOption) error
	Authorize() error
}

type GsyncServiceContract interface {
	Pull() error
	Push() error
	Sync() error
}
