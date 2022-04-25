package domain

type SynchronizableStoreContract interface {
	Pull(option SyncOption) error
	Push(option SyncOption) error
	Authorize() error
}

type GsyncServiceContract interface {
	Pull(option ...SyncOption) error
	Push(option ...SyncOption) error
	Sync(option ...SyncOption) error
	Authorize() error
}
