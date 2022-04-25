package domain

type GsyncStoreContract interface {
	Pull() error
	Push() error
	Sync() error
	Authorize() error
}

type GsyncServiceContract interface {
	Pull() error
	Push() error
	Sync() error
	authorize() error
}
