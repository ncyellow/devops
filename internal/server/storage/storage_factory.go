// Package storage содержит разные имплементации сохранения метрик
// В никуда, на дис в файл, либо в БД
// Стандартный вариант использования
// repo := repository.NewRepository(s.Conf.GeneralCfg())
//	saver, err := storage.CreateStorage(s.Conf, repo)
package storage

import (
	"github.com/ncyellow/devops/internal/repository"
	"github.com/ncyellow/devops/internal/server/config"
)

// CreateStorage фабричная функция которая по настройкам возвращает имплементацию хранилища
// либо Postgres, либо файл, либо фейковый пустой вариант.
func CreateStorage(conf *config.Config, repo repository.Repository) (PersistentStorage, error) {
	if conf.DatabaseConn != "" {
		return NewPgStorage(conf, repo)
	} else if conf.StoreFile != "" {
		return NewFileStorage(conf, repo)
	}
	return NewFakeStorage()
}
