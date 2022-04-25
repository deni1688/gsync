package domain

type gsyncService struct {
	store SynchronizableStorageContract
}

func NewGsyncService(store SynchronizableStorageContract) GsyncServiceContract {
	return &gsyncService{store: store}
}

func (g gsyncService) Pull() error {
	//TODO implement me
	panic("implement me")
}

func (g gsyncService) Push() error {
	//TODO implement me
	panic("implement me")
}

func (g gsyncService) Sync() error {
	//TODO implement me
	panic("implement me")
}
