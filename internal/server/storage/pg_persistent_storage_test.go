package storage

import (
	"errors"
	"testing"

	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/ncyellow/devops/internal/repository"
	"github.com/ncyellow/devops/internal/server/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PgStorageSuite struct {
	suite.Suite
	saver    PersistentStorage
	mockPool *pgxpoolmock.MockPgxPool
	repo     repository.Repository
}

func TestPgStorageSuite(t *testing.T) {
	suite.Run(t, new(PgStorageSuite))
}

// SetupSuite перед началом теста стартуем новый сервер httptest.Server делаем так, чтобы тестировать каждый
// handler отдельно и не сливать все тесты в один
func (suite *PgStorageSuite) SetupTest() {

	conf := config.Config{}
	repo := repository.NewRepository(conf.GeneralCfg())
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	// given
	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

	// При инициализации у нас идут create скрипты
	tag := []byte("create")
	mockPool.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any()).Return(tag, nil).AnyTimes()

	saver := &PgPersistentStorage{conf: &conf, pool: mockPool, repo: repo}
	saver.init()
	suite.saver = saver
	suite.mockPool = mockPool
	suite.repo = repo
}

// TearDownSuite после теста отключаем сервер
func (suite *PgStorageSuite) TearDownTest() {
	// Проверяем что при Close закрывается пул
	// Очистка метрик чтобы они не вызывали сохранение - там очень сложная логика
	suite.repo.Clear()
	suite.mockPool.EXPECT().Close()
	suite.saver.Close()
}

func (suite *PgStorageSuite) TestPing() {
	// При пинге проверяем что у нас у нас идет "select 1" запрос
	tag := []byte("select")
	suite.mockPool.EXPECT().Exec(gomock.Any(), "select 1", gomock.Any()).Return(tag, nil).AnyTimes()
	suite.saver.Ping()
}

func (suite *PgStorageSuite) TestLoad() {

	counterColumns := []string{"metric_name", "value"}
	pgxRowsCounter := pgxpoolmock.NewRows(counterColumns).
		AddRow("testCounter", int64(100)).
		AddRow("minCounter", int64(10)).
		ToPgxRows()

	suite.mockPool.EXPECT().Query(gomock.Any(), `select "metric_name", "value" FROM "counters"`, gomock.Any()).
		Return(pgxRowsCounter, nil)

	gaugeColumns := []string{"metric_name", "value"}
	pgxRowsGauge := pgxpoolmock.NewRows(gaugeColumns).
		AddRow("testGauge", 110.5).
		AddRow("maxGauge", 11.5).
		ToPgxRows()

	suite.mockPool.EXPECT().Query(gomock.Any(), `select "metric_name", "value" FROM "gauges"`, gomock.Any()).
		Return(pgxRowsGauge, nil)
	err := suite.saver.Load()
	assert.NoError(suite.T(), err)

	// Должно быть 4 метрики в репозитории 2 gauge + 2 counter
	assert.Equal(suite.T(), len(suite.repo.ToMetrics()), 4)

	// Убираем метрики чтобы каждый тест начинался с чистыми метриками
	suite.repo.Clear()
}

func (suite *PgStorageSuite) TestCounterFailedLoad() {

	counterColumns := []string{"metric_name", "value"}
	pgxRowsCounter := pgxpoolmock.NewRows(counterColumns).
		AddRow("testCounter", 100).
		AddRow("minCounter", 10).
		ToPgxRows()

	suite.mockPool.EXPECT().Query(gomock.Any(), `select "metric_name", "value" FROM "counters"`, gomock.Any()).
		Return(pgxRowsCounter, nil)

	gaugeColumns := []string{"metric_name", "value"}
	pgxRowsGauge := pgxpoolmock.NewRows(gaugeColumns).
		AddRow("testGauge", 110.5).
		AddRow("maxGauge", 11.5).
		ToPgxRows()

	suite.mockPool.EXPECT().Query(gomock.Any(), `select "metric_name", "value" FROM "gauges"`, gomock.Any()).
		Return(pgxRowsGauge, nil)
	err := suite.saver.Load()

	// Будет ошибка, так как тип метрик counter Не соответствует должен быть int64 вместо int
	assert.Error(suite.T(), err)

	// список метрик пустой так как выпала ошибка
	assert.Equal(suite.T(), len(suite.repo.ToMetrics()), 0)

	// Убираем метрики чтобы каждый тест начинался с чистыми метриками
	suite.repo.Clear()
}

func (suite *PgStorageSuite) TestGaugeFailedLoad() {

	counterColumns := []string{"metric_name", "value"}
	pgxRowsCounter := pgxpoolmock.NewRows(counterColumns).
		AddRow("testCounter", int64(100)).
		AddRow("minCounter", int64(10)).
		ToPgxRows()

	suite.mockPool.EXPECT().Query(gomock.Any(), `select "metric_name", "value" FROM "counters"`, gomock.Any()).
		Return(pgxRowsCounter, nil)

	gaugeColumns := []string{"metric_name", "value"}
	pgxRowsGauge := pgxpoolmock.NewRows(gaugeColumns).
		AddRow("testGauge", "test").
		AddRow("maxGauge", "test").
		ToPgxRows()

	suite.mockPool.EXPECT().Query(gomock.Any(), `select "metric_name", "value" FROM "gauges"`, gomock.Any()).
		Return(pgxRowsGauge, nil)
	err := suite.saver.Load()

	// Будет ошибка, так как тип метрик gauge Не соответствует должен быть float64 вместо int
	assert.Error(suite.T(), err)

	// список метрик пустой так как выпала ошибка
	assert.Equal(suite.T(), len(suite.repo.ToMetrics()), 0)

	// Убираем метрики чтобы каждый тест начинался с чистыми метриками
	suite.repo.Clear()
}

func (suite *PgStorageSuite) TestGaugeQueryFailedLoad() {

	counterColumns := []string{"metric_name", "value"}
	pgxRowsCounter := pgxpoolmock.NewRows(counterColumns).
		AddRow("testCounter", int64(100)).
		AddRow("minCounter", int64(10)).
		ToPgxRows()

	suite.mockPool.EXPECT().Query(gomock.Any(), `select "metric_name", "value" FROM "counters"`, gomock.Any()).
		Return(pgxRowsCounter, nil)

	suite.mockPool.EXPECT().Query(gomock.Any(), `select "metric_name", "value" FROM "gauges"`, gomock.Any()).
		Return(nil, errors.New("some error"))
	err := suite.saver.Load()

	// Будет ошибка, так как упал запрос за метриками gauge
	assert.Error(suite.T(), err)

	// список метрик пустой так как выпала ошибка
	assert.Equal(suite.T(), len(suite.repo.ToMetrics()), 0)

	// Убираем метрики чтобы каждый тест начинался с чистыми метриками
	suite.repo.Clear()
}

func (suite *PgStorageSuite) TestCounterQueryFailedLoad() {

	suite.mockPool.EXPECT().Query(gomock.Any(), `select "metric_name", "value" FROM "counters"`, gomock.Any()).
		Return(nil, errors.New("some error"))

	err := suite.saver.Load()

	// Будет ошибка, так как упал запрос за метриками gauge
	assert.Error(suite.T(), err)

	// список метрик пустой так как выпала ошибка
	assert.Equal(suite.T(), len(suite.repo.ToMetrics()), 0)

	// Убираем метрики чтобы каждый тест начинался с чистыми метриками
	suite.repo.Clear()
}
