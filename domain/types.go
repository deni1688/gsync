package domain

type SyncPriority int

const (
	SyncPriorityRemote SyncPriority = iota
	SyncPriorityLocal
)

type SyncOption struct {
	FilePath string
	Priority SyncPriority
}

func (s SyncOption) GetPriority() string {
	switch s.Priority {
	case SyncPriorityRemote:
		return "remote"
	case SyncPriorityLocal:
		return "local"
	default:
		return "unknown"
	}
}
