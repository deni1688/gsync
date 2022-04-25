package domain

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
