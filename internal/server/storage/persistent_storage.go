package storage

// PersistentStorage интерфейс хранилища, для загрузки и сохранения данных
type PersistentStorage interface {
	Save() error
	Load() error
	Ping() error
	Close()
}
