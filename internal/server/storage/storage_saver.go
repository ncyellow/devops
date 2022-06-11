package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ncyellow/devops/internal/server/config"
)

func SaveToFile(fileName string, repo Repository) {
	//! Если файл не задан, ок ничего не делаем
	if fileName == "" {
		return
	}

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	defer file.Close()
	if err != nil {
		fmt.Printf("can't open file %s", fileName)
		return
	}
	encoder := json.NewEncoder(file)
	encoder.Encode(&repo)
}

func RestoreFromFile(fileName string, repo Repository) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0777)
	defer file.Close()
	if err != nil {
		fmt.Printf("can't open file %s", fileName)
		return
	}
	decoder := json.NewDecoder(file)
	decoder.Decode(&repo)
}

func RunStorageSaver(config config.Config, repo Repository) {
	if config.StoreInterval == 0 {
		//! Не нужно сбрасывать на диск если StoreInterval == 0
		return
	}

	tickerStore := time.NewTicker(config.StoreInterval)
	defer tickerStore.Stop()

	for {
		select {
		case <-tickerStore.C:
			//! сбрасываем на диск
			SaveToFile(config.StoreFile, repo)
		}
	}
}
