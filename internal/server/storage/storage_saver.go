package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ncyellow/devops/internal/server/config"
)

type NullSaver struct {
}

func (m *NullSaver) Close(repo Repository) {
}

func (m *NullSaver) Load(repo Repository) error {
	return nil
}

func (m *NullSaver) Save(repo Repository) error {
	return nil
}

func NewNullSaver() (Saver, error) {
	return &NullSaver{}, nil
}

type MemoryStorageSaver struct {
	conf *config.Config
}

func NewMemorySaver(conf *config.Config) (Saver, error) {
	saver := MemoryStorageSaver{conf: conf}
	return &saver, nil
}

func CreateSaver(conf *config.Config) (Saver, error) {
	if conf.DatabaseConn != "" {
		return NewSaver(conf)
	} else if conf.StoreFile != "" {
		return NewMemorySaver(conf)
	}
	return NewNullSaver()
}

func (m *MemoryStorageSaver) Close(repo Repository) {
}

func (m *MemoryStorageSaver) Load(repo Repository) error {
	if m.conf.Restore {
		RestoreFromFile(m.conf.StoreFile, repo)
	}
	return nil
}

func (m *MemoryStorageSaver) Save(repo Repository) error {
	SaveToFile(m.conf.StoreFile, repo)
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

// RunStorageSaver запускает сохранение данных repo по таймеру в файл
func RunStorageSaver(config *config.Config, repo Repository) {
	if config.StoreInterval == 0 {
		//! Не нужно сбрасывать на диск если StoreInterval == 0
		return
	}

	tickerStore := time.NewTicker(config.StoreInterval)
	defer tickerStore.Stop()

	for {
		<-tickerStore.C
		//! сбрасываем на диск
		SaveToFile(config.StoreFile, repo)
	}
}
