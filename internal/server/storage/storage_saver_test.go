package storage

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveRestoreFromFile(t *testing.T) {

	// Создаем репозиторий, который будем тестировать
	repo := NewRepository()
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
	SaveToFile(fileName, repo)

	// читаем из файла и сравниваем метрики
	newRepo := NewRepository()
	RestoreFromFile(fileName, newRepo)

	// Так как перегружен Stringer, который возвращает нам html они должны быть одинаковые
	assert.Equal(t, repo.String(), newRepo.String())

}