package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ncyellow/devops/internal/server/config"
)

type NullSaver struct {
}

func (m *NullSaver) Close() {
}

func (m *NullSaver) Load() error {
	return nil
}

func (m *NullSaver) Save() error {
	return nil
}

func NewNullSaver() (PersistentStorage, error) {
	return &NullSaver{}, nil
}

type MemoryStorageSaver struct {
	conf *config.Config
	repo Repository
}

func NewMemorySaver(conf *config.Config, repo Repository) (PersistentStorage, error) {
	saver := MemoryStorageSaver{conf: conf, repo: repo}
	return &saver, nil
}

func CreateSaver(conf *config.Config, repo Repository) (PersistentStorage, error) {
	if conf.DatabaseConn != "" {
		return NewSaver(conf, repo)
	} else if conf.StoreFile != "" {
		return NewMemorySaver(conf, repo)
	}
	return NewNullSaver()
}

func (m *MemoryStorageSaver) Close() {
}

func (m *MemoryStorageSaver) Load() error {
	if m.conf.Restore {
		RestoreFromFile(m.conf.StoreFile, m.repo)
	}
	return nil
}

func (m *MemoryStorageSaver) Save() error {
	SaveToFile(m.conf.StoreFile, m.repo)
	return nil
}

// SaveToFile сохраняет данные repo в файл с именем fileName
func SaveToFile(fileName string, repo Repository) {
	//! Если файл не задан, ок ничего не делаем
	if fileName == "" {
		return
	}

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		fmt.Printf("can't open file %s", fileName)
		return
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.Encode(&repo)
}

// RestoreFromFile загружает данные в repo из файла с именем fileName
func RestoreFromFile(fileName string, repo Repository) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0777)
	if err != nil {
		fmt.Printf("can't open file %s", fileName)
		return
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	decoder.Decode(&repo)
}
