package domain

func (gs gsyncService) Sync(dir SyncTarget) error {
	if err := gs.Push(dir); err != nil {
		return err
	}

	if err := gs.Pull(dir); err != nil {
		return err
	}

	return nil
}
