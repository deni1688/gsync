package domain

type gsyncService struct {
	store SynchronizableStoreContract
}

func NewGsyncService(store SynchronizableStoreContract) GsyncServiceContract {
	return &gsyncService{store}
}

func (g gsyncService) Pull(options ...SyncOption) error {
	var err error

	for _, option := range options {
		err = g.store.Pull(option)
	}

	return err
}

func (g gsyncService) Push(options ...SyncOption) error {
	var err error

	for _, option := range options {
		err = g.store.Push(option)
	}

	return err
}

func (g gsyncService) Sync(options ...SyncOption) error {
	err := g.Pull(options...)
	err = g.Push(options...)

	return err
}

func (g gsyncService) Authorize() error {
	return g.store.Authorize()
}
