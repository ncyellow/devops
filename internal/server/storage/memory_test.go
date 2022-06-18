// Тесты для MapRepository
package storage

import (
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MapRepository", func() {
	Context("Если выполняется сериализация в json MapRepository", func() {
		It("json должен быть корректным и содержать все значения из MapRepository", func() {
			repo := NewRepository()

			err := repo.UpdateGauge("testGaugeMetric", 100)
			Expect(err).Should(BeNil())

			err = repo.UpdateCounter("testCounterMetric", 120)
			Expect(err).Should(BeNil())

			jsRepo, err := json.Marshal(repo)
			Expect(err).Should(BeNil())

			Expect(jsRepo).Should(MatchJSON(`[{"id":"testGaugeMetric","type":"gauge","value":100},{"id":"testCounterMetric","type":"counter","delta":120}]`))
		})
	})

	Context("Если выполняется десериализация json в MapRepository", func() {
		It("Разобраные метрики должны соответствовать исходным значениям", func() {
			repo := NewRepository()
			data := []byte(`[{"id":"testGaugeMetric","type":"gauge","value":100},{"id":"testCounterMetric","type":"counter","delta":120}]`)

			err := json.Unmarshal(data, &repo)
			Expect(err).Should(BeNil())

			val, ok := repo.Gauge("testGaugeMetric")
			Expect(ok).Should(Equal(true))
			Expect(val).Should(Equal(100.0))

			delta, ok := repo.Counter("testCounterMetric")
			Expect(ok).Should(Equal(true))
			Expect(delta).Should(Equal(int64(120)))
		})
	})

	Context("Если передан некорректный json", func() {
		It("Unmarshal завершаться с ошибкой", func() {
			brokenRepo := NewRepository()
			brokenData := []byte(`{"name": "Joe", "age": null, }`)

			err := json.Unmarshal(brokenData, &brokenRepo)
			Expect(err).ShouldNot(BeNil())
		})
	})

	Context("Выполняем проверку работы gauge метрик", func() {
		It("Проверка простого чтения и записи через API Metric", func() {
			repo := NewRepository()

			var updateValue float64 = 100

			// обновление
			err := repo.UpdateMetric(Metrics{
				ID:    "testGaugeMetric",
				MType: Gauge,
				Value: &updateValue,
			})
			Expect(err).Should(BeNil())

			// чтение
			val, ok := repo.Metric("testGaugeMetric", Gauge)
			Expect(ok).Should(Equal(true))
			Expect(*val.Value).Should(Equal(updateValue))
		})

		It("После повторного обновления gauge должно быть выставлено последнее значение", func() {
			repo := NewRepository()
			// обновление
			var updateValue float64 = 100

			// обновление
			err := repo.UpdateMetric(Metrics{
				ID:    "testGaugeMetric",
				MType: Gauge,
				Value: &updateValue,
			})
			Expect(err).Should(BeNil())

			// обновляем повторно
			updateValue = 300
			err = repo.UpdateMetric(Metrics{
				ID:    "testGaugeMetric",
				MType: Gauge,
				Value: &updateValue,
			})
			Expect(err).Should(BeNil())

			// проверяем что старое значение перезаписалось
			val, ok := repo.Metric("testGaugeMetric", Gauge)
			Expect(*val.Value).Should(Equal(updateValue))
			Expect(ok).Should(Equal(true))
		})

		It("Для неизвестного типа возвращается ошибка метрик", func() {
			repo := NewRepository()
			updateValue := 100.0
			err := repo.UpdateMetric(Metrics{
				ID:    "testMetric",
				MType: "unknownType",
				Value: &updateValue,
			})
			Expect(err).ShouldNot(BeNil())
		})

		It("Чтение неизвестной метрики типа gauge возвращает Metrics{}, false", func() {
			repo := NewRepository()
			// Проверка чтения неизвестной метрики тика Gauge
			val, ok := repo.Metric("unknownMetricGauge", Gauge)
			Expect(ok).Should(BeFalse())
			Expect(val).Should(Equal(Metrics{}))
		})
	})

	Context("Выполняем проверку работы counter метрик", func() {

		It("Чтение неизвестной метрики типа counter возвращает Metrics{}, false", func() {
			repo := NewRepository()
			// Проверка чтения неизвестной метрики тика Gauge
			val, ok := repo.Metric("unknownMetricCounter", Counter)
			Expect(ok).Should(BeFalse())
			Expect(val).Should(Equal(Metrics{}))
		})

		It("Проверка стандартного чтения записи counter метрик", func() {
			repo := NewRepository()

			// обновление
			err := repo.UpdateCounter("testCounter", 100)
			Expect(err).Should(BeNil())

			// чтение
			val, ok := repo.Counter("testCounter")
			Expect(ok).Should(Equal(true))
			Expect(val).Should(Equal(int64(100)))
		})
		It("Проверка что два обновления метрики увеличивает счетчик counter", func() {
			repo := NewRepository()

			// обновление
			err := repo.UpdateCounter("testCounter", 100)
			Expect(err).Should(BeNil())

			// обновление
			err = repo.UpdateCounter("testCounter", 100)
			Expect(err).Should(BeNil())

			// чтение
			val, ok := repo.Counter("testCounter")
			Expect(ok).Should(Equal(true))
			Expect(val).Should(Equal(int64(200)))
		})

		It("Проверка чтения неизвестной метрики counter", func() {
			repo := NewRepository()
			// Проверка чтения неизвестного значения
			_, ok := repo.Counter("unknownCounter")
			Expect(ok).Should(Equal(false))
		})

		It("Проверка стандартного чтения записи counter метрик через API Metrics", func() {
			repo := NewRepository()

			var updateValue int64 = 100

			// обновление
			err := repo.UpdateMetric(Metrics{
				ID:    "testCounterMetric",
				MType: Counter,
				Delta: &updateValue,
			})
			Expect(err).Should(BeNil())

			// чтение
			val, ok := repo.Metric("testCounterMetric", Counter)
			Expect(ok).Should(Equal(true))
			Expect(*val.Delta).Should(Equal(updateValue))
		})
	})

	Context("Проверяем интерфейс Stringer", func() {
		It("Должен быть корректный html наличии метрик gauge и counter", func() {
			repo := NewRepository()

			// обновление
			err := repo.UpdateGauge("testGauge", 100.0)
			Expect(err).Should(BeNil())

			err = repo.UpdateCounter("testCounter", 100)
			Expect(err).Should(BeNil())

			correctHTML := `
	<html>
	<body>
	<h1>All metrics</h1>
	<h3>gauges</h3>
	<ul>
	  <li>testGauge : 100.000</li>

	</ul>
	<h3>counters</h3>
	<ul>
	  <li>testCounter : 100</li>

	</ul>
	</body>
	</html>`

			Expect(repo.String()).Should(Equal(correctHTML))
		})
	})
})
