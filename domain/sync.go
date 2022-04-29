package domain

func (g syncService) Sync(syncFile SyncFile) error {
	var err error

	err = g.Push(syncFile)
	err = g.Pull(syncFile)

	return err
}
