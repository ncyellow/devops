package storage

import (
	"testing"

	"github.com/ncyellow/devops/internal/repository"
	"github.com/ncyellow/devops/internal/server/config"
	"github.com/stretchr/testify/assert"
)

func TestCreateStorage(t *testing.T) {
	conf := config.Config{}
	repo := repository.NewRepository(conf.GeneralCfg())

	storage, err := CreateStorage(&config.Config{}, repo)
	assert.NoError(t, err)
	assert.IsType(t, &FakeStorage{}, storage)

	storage, err = CreateStorage(&config.Config{
		StoreFile: "/tmp/testdata",
	}, repo)
	assert.NoError(t, err)
	assert.IsType(t, &FileStorageSaver{}, storage)

	// тип посгреса мы создать не можем так как у нас нет базы в юнит тестах - потому будет ошибка подключения
	storage, err = CreateStorage(&config.Config{
		DatabaseConn: "superconn",
	}, repo)
	assert.Error(t, err)
	assert.Nil(t, storage)
}
