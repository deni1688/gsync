package domain

type gsyncService struct {
	store GsyncStoreContract
}

func NewGsyncService(store GsyncStoreContract) GsyncServiceContract {
	return &gsyncService{store: store}
}

func (g gsyncService) authorize() error {
	//TODO implement me
	panic("implement me")
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
