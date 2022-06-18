package storage

import (
	"encoding/json"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestStorageSaver(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Storage Saver Suite")
}

var _ = Describe("Storage Saver", func() {

	repo := NewRepository()
	var fileName string

	BeforeSuite(func() {
		data := []byte(`[{"id":"testGaugeMetric","type":"gauge","value":100},{"id":"testCounterMetric","type":"counter","delta":120}]`)
		err := json.Unmarshal(data, &repo)
		Expect(err).Should(BeNil())

		// Создаем временный файл
		file, err := os.CreateTemp(os.TempDir(), "restore*")
		Expect(err).Should(BeNil())

		fileName = file.Name()
		// Закрываем так как открывать мы будем его сами
		file.Close()
	})

	AfterSuite(func() {
		// Удаляем тестовый временный файл
		os.Remove(fileName)
	})

	Context("Если при восстановении передан не существующий файл", func() {
		It("Метод должен завершаться с ошибкой", func() {
			newRepo := NewRepository()
			Expect(RestoreFromFile("", newRepo)).ShouldNot(BeNil())
		})
	})

	Context("Если при сохранении передан левый файл", func() {
		It("Метод должен завершаться без ошибки", func() {
			Expect(SaveToFile("", repo)).Should(BeNil())
		})
	})

	It("После сохранения восстановления данные должны быть идентичны", func() {
		Expect(SaveToFile(fileName, repo)).Should(BeNil())

		newRepo := NewRepository()
		Expect(RestoreFromFile(fileName, newRepo)).Should(BeNil())
		Expect(repo.String()).Should(Equal(newRepo.String()))
	})
})
