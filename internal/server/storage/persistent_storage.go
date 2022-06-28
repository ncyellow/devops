package storage

type PersistentStorage interface {
	Save() error
	Load() error
	Close()
}
