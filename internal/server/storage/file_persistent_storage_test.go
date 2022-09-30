package storage

import (
	"context"
	"encoding/json"
	"github.com/ncyellow/devops/internal/repository"
	"github.com/ncyellow/devops/internal/server/config"

	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveRestoreFromFile(t *testing.T) {

	// Создаем репозиторий, который будем тестировать
	cfg := config.Config{}
	repo := repository.NewRepository(cfg.GeneralCfg())

	data := []byte(`[{"id":"testGaugeMetric","type":"gauge","value":100},{"id":"testCounterMetric","type":"counter","delta":120}]`)
	err := json.Unmarshal(data, &repo)
	assert.NoError(t, err)

	// Создаем временный файл
	file, err := os.CreateTemp(os.TempDir(), "restore*")
	assert.NoError(t, err)
	fileName := file.Name()
	// Закрываем так как открывать мы будем его сами
	file.Close()

	// Удаляем файл в конце теста
	defer os.Remove(fileName)

	// Сохраняем в файл
	err = SaveToFile(fileName, repo)
	assert.NoError(t, err)

	// читаем из файла и сравниваем метрики
	newRepo := repository.NewRepository(cfg.GeneralCfg())
	RestoreFromFile(fileName, newRepo)

	// Так как перегружен Stringer, который возвращает нам html они должны быть одинаковые
	// Второй вариант сравнить их json представление
	assert.Equal(t, repository.RenderHTML(repo.ToMetrics()), repository.RenderHTML(newRepo.ToMetrics()))

}

func TestSaveToFile(t *testing.T) {
	conf := config.Config{}
	repo := repository.NewRepository(conf.GeneralCfg())
	assert.NoError(t, SaveToFile("file doesn't not exists", repo))
}

func TestNewFileStorage(t *testing.T) {
	conf := config.Config{}
	repo := repository.NewRepository(conf.GeneralCfg())
	storage, err := NewFileStorage(&conf, repo)
	assert.NoError(t, err)
	defer storage.Close()
	assert.Nil(t, storage.Ping())
	assert.Nil(t, storage.Load())
	assert.Nil(t, storage.Save(context.Background()))
}
