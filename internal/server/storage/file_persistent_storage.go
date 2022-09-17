package storage

import (
	"context"
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/ncyellow/devops/internal/repository"
	"github.com/ncyellow/devops/internal/server/config"
)

// FileStorageSaver реализация хранилища на основе файла, реализует интерфейс PersistentStorage
type FileStorageSaver struct {
	conf *config.Config
	repo repository.Repository
}

// NewFileStorage конструктор хранилища на основе файла, явно не используется, только через фабрику
func NewFileStorage(conf *config.Config, repo repository.Repository) (*FileStorageSaver, error) {
	saver := FileStorageSaver{conf: conf, repo: repo}
	return &saver, nil
}

func (m *FileStorageSaver) Ping() error {
	return nil
}

func (m *FileStorageSaver) Close() {
}

func (m *FileStorageSaver) Load() error {
	if m.conf.Restore {
		RestoreFromFile(m.conf.StoreFile, m.repo)
	}
	return nil
}

func (m *FileStorageSaver) Save(context.Context) error {
	return SaveToFile(m.conf.StoreFile, m.repo)
}

// SaveToFile сохраняет данные repo в файл с именем fileName
func SaveToFile(fileName string, repo repository.Repository) error {
	//! Если файл не задан, ок ничего не делаем
	if fileName == "" {
		return nil
	}

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Info().Msgf("can't open file %s", fileName)
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.Encode(&repo)
	return nil
}

// RestoreFromFile загружает данные в repo из файла с именем fileName
func RestoreFromFile(fileName string, repo repository.Repository) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0777)
	if err != nil {
		log.Info().Msgf("can't open file %s", fileName)
		return
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	decoder.Decode(&repo)
}
