package storage

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/ncyellow/devops/internal/server/config"
	"github.com/ncyellow/devops/internal/server/repository"
)

type FileStorageSaver struct {
	conf *config.Config
	repo repository.Repository
}

func NewFileStorage(conf *config.Config, repo repository.Repository) (PersistentStorage, error) {
	saver := FileStorageSaver{conf: conf, repo: repo}
	return &saver, nil
}

func (m *FileStorageSaver) Close() {
}

func (m *FileStorageSaver) Load() error {
	if m.conf.Restore {
		RestoreFromFile(m.conf.StoreFile, m.repo)
	}
	return nil
}

func (m *FileStorageSaver) Save() error {
	SaveToFile(m.conf.StoreFile, m.repo)
	return nil
}

// SaveToFile сохраняет данные repo в файл с именем fileName
func SaveToFile(fileName string, repo repository.Repository) {
	//! Если файл не задан, ок ничего не делаем
	if fileName == "" {
		return
	}

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Info().Msgf("can't open file %s", fileName)
		return
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.Encode(&repo)
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
