package storage

import (
	"github.com/ncyellow/devops/internal/server/config"
	"github.com/ncyellow/devops/internal/server/repository"
)

func CreateStorage(conf *config.Config, repo repository.Repository) (PersistentStorage, error) {
	if conf.DatabaseConn != "" {
		return NewSaver(conf, repo)
	} else if conf.StoreFile != "" {
		return NewMemorySaver(conf, repo)
	}
	return NewFakeStorage()
}
