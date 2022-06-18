package storage

import (
	"encoding/json"
	"os"
	"time"

	"github.com/ncyellow/devops/internal/server/config"
)

// SaveToFile сохраняет данные repo в файл с именем fileName
func SaveToFile(fileName string, repo Repository) error {
	//! Если файл не задан, ок ничего не делаем
	if fileName == "" {
		return nil
	}

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	return encoder.Encode(&repo)
}

// RestoreFromFile загружает данные в repo из файла с именем fileName
func RestoreFromFile(fileName string, repo Repository) error {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0777)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	return decoder.Decode(&repo)
}

// RunStorageSaver запускает сохранение данных repo по таймеру в файл
func RunStorageSaver(config config.Config, repo Repository) {
	if config.StoreInterval == 0 {
		//! Не нужно сбрасывать на диск если StoreInterval == 0
		return
	}

	tickerStore := time.NewTicker(config.StoreInterval)
	defer tickerStore.Stop()

	for {
		<-tickerStore.C
		//! сбрасываем на диск
		_ = SaveToFile(config.StoreFile, repo)
	}
}
