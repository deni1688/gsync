package domain

func (g syncService) Sync(dir SyncFile) error {
	if err := g.Push(dir); err != nil {
		return err
	}

	if err := g.Pull(dir); err != nil {
		return err
	}

	return nil
}
