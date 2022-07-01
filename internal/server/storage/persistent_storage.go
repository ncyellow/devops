package storage

import "context"

// PersistentStorage интерфейс хранилища, для загрузки и сохранения данных
type PersistentStorage interface {
	// Save сохранение данных в хранилище
	Save(ctx context.Context) error
	// Load загрузка данных из хранилища
	Load() error
	// Ping проверка доступности хранилища
	Ping() error
	// Close вызывается при окончании работы для закрытия коннектов и закрытия файлов
	Close()
}
