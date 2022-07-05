package storage

import (
	"github.com/ncyellow/devops/internal/server/config"
	"github.com/ncyellow/devops/internal/server/repository"
)

// CreateStorage фабричная функция которая по настройкам возвращает имплементацию хранилища
func CreateStorage(conf *config.Config, repo repository.Repository) (PersistentStorage, error) {
	if conf.DatabaseConn != "" {
		return NewPgStorage(conf, repo)
	} else if conf.StoreFile != "" {
		return NewFileStorage(conf, repo)
	}
	return NewFakeStorage()
}
