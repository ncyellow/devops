// Package handlers содержит роутинг и все обработчики запросов сервера
package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/ncyellow/devops/internal/hash"
	"github.com/ncyellow/devops/internal/repository"
	"github.com/ncyellow/devops/internal/server/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ncyellow/devops/internal/server/middlewares"
	"github.com/ncyellow/devops/internal/server/storage"
)

var (
	AnswerOK = []byte("ok")
)

// @Title DevOPS API
// @Description Сервис сбора метрик типов Counter, Gauge
// @Version 1.0

// @Contact.email ncyellow@yandex.ru

// @Tag.name Info
// @Tag.description "Группа запросов на получение состояние сервера и метрик"

// @Tag.name Storage
// @Tag.description "Группа запросов на изменение метрик"

// Handler структура данных для работы с роутингом
type Handler struct {
	*chi.Mux
	conf   *config.Config
	repo   repository.Repository
	pStore storage.PersistentStorage
}

// NewRouter создает chi.NewRouter и описывает маршрутизацию
func NewRouter(repo repository.Repository, conf *config.Config, pStore storage.PersistentStorage) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middlewares.EncoderGZIP)
	r.Mount("/debug", middleware.Profiler())

	handler := &Handler{
		Mux:    r,
		conf:   conf,
		repo:   repo,
		pStore: pStore,
	}
	handler.Get("/", handler.List())
	handler.Get("/value/{metricType}/{metricName}", handler.Value())
	r.Post("/update/{metricType}/{metricName}/{metricValue}", handler.Update())
	r.Post("/update/", handler.UpdateJSON())
	r.Post("/value/", handler.ValueJSON())
	r.Post("/updates/", handler.UpdateListJSON())
	r.Get("/ping", handler.Ping())

	return handler
}

// List возвращает html произвольного формата со всеми метриками сервера
// @Tags Info
// @Summary Возвращает html со списком метрик
// @Description Просто генерим рандомного формата html с метриками
// @ID infoList
// @Produce plain
// @Success 200 {string} string "html с метриками"
// @Router / [get]
func (h *Handler) List() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "text/html")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(repository.RenderHTML(h.repo.ToMetrics())))
	}
}

// Value возвращает значение конкретной метрики через GET
// @Tags Info
// @Summary Возвращает состояние метрики текстом
// @Description на вход rest url на выход plain значение
// @ID infoValue
// @Produce plain
// @Param metricType path string true "Metric type"
// @Param metricName path string true "Metric name"
// @Success 200 {string} string "Значение метрики к примеру - 10.2"
// @Failure 404 {string} string "not found"
// @Router /value/{metricType}/{metricName} [get]
func (h *Handler) Value() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")
		switch metricType {
		case repository.Gauge:
			val, ok := h.repo.Gauge(metricName)
			if ok {
				rw.Write([]byte(fmt.Sprintf("%.03f", val)))
				return
			} else {
				rw.WriteHeader(http.StatusNotFound)
				rw.Write([]byte("not found"))
				return
			}
		case repository.Counter:
			val, ok := h.repo.Counter(metricName)
			if ok {
				rw.Write([]byte(fmt.Sprintf("%d", val)))
				return
			} else {
				rw.WriteHeader(http.StatusNotFound)
				rw.Write([]byte("not found"))
				return
			}
		default:
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte("not found"))
		}
	}
}

// Update обновляет значение конкретной метрики в rest формате
// @Tags Storage
// @Summary обновляем состояние метрики через rest api
// @Description на вход rest url на выход plain ок если все хорошо
// @ID storageValue
// @Produce plain
// @Param metricType path string true "Metric type"
// @Param metricName path string true "Metric name"
// @Param metricValue path string true "Metric value"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "incorrect metric value"
// @Failure 500 {string} string "incorrect metric name"
// @Failure 501 {string} string "incorrect metric type"
// @Router /update/{metricType}/{metricName}/{metricValue} [post]
func (h *Handler) Update() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		//! Метод только post
		if r.Method != http.MethodPost {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("only post method support"))
			return
		}

		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")
		metricValue := chi.URLParam(r, "metricValue")

		switch metricType {
		case repository.Gauge:
			value, err := strconv.ParseFloat(metricValue, 64)
			//! Второй параметр обязательно кастится в float64
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write([]byte("incorrect metric value"))
				return
			}
			err = h.repo.UpdateGauge(metricName, value)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write([]byte("incorrect metric name "))
				return
			}
		case repository.Counter:
			value, err := strconv.ParseInt(metricValue, 10, 64)
			//! Второй параметр обязательно кастится в int64
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write([]byte("incorrect metric value"))
				return
			}
			err = h.repo.UpdateCounter(metricName, value)
			//! Сейчас проблема только одна - ошибка при кривом имени метрики
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write([]byte("incorrect metric name"))
				return
			}
		default:
			rw.WriteHeader(http.StatusNotImplemented)
			rw.Write([]byte("incorrect metric type"))
			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write(AnswerOK)
	}
}

// UpdateJSON возвращает значение конкретной метрики, но запрос приходит в json body
// @Tags Storage
// @Summary обновляем состояние метрики но уже через json body
// @Description важный момент что запрос на состояние метрики должен быть подписан корректно иначе, отлуп
// @ID storageUpdateJSON
// @Accept json
// @Produce plain
// @Param metric_data body Metrics true "Metric object"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "incorrect metric sign"
// @Failure 500 {string} string "incorrect metric type, content type not support, invalid deserialization"
// @Router /update/ [post]
func (h *Handler) UpdateJSON() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("content type not support"))
			return
		}
		reqBody, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Read data problem"))
			return
		}
		metric := repository.Metrics{}
		err = json.Unmarshal(reqBody, &metric)

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("invalid deserialization"))
			return
		}

		encodeFunc := hash.CreateEncodeFunc(h.conf.SecretKey)
		ok := hash.CheckSign(h.conf.SecretKey, metric.Hash, metric.CalcHash(encodeFunc))
		if !ok {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("incorrect metric sign"))
			return
		}

		err = h.repo.UpdateMetric(metric)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("incorrect metric type"))
			return
		}
		h.pStore.Save(r.Context())

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("ok"))
	}
}

// UpdateListJSON обновляет значение всех метрик переданных в json body
// @Tags Storage
// @Summary обновляем состояние всех метрик переданных  в json
// @Description обязательность подписи как и UpdateJSON остается
// @ID storageUpdateListJSON
// @Accept json
// @Produce plain
// @Param metric_data body []Metrics true "Metrics list object"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "incorrect metric sign"
// @Failure 500 {string} string "incorrect metric type, content type not support, invalid deserialization"
// @Router /updates/ [post]
func (h *Handler) UpdateListJSON() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Content-Type") != "application/json" {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("content type not support"))
			return
		}
		reqBody, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Read data problem"))
			return
		}
		var metrics []repository.Metrics
		err = json.Unmarshal(reqBody, &metrics)

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("invalid deserialization"))
			return
		}

		//! Проверяем подписи - если есть криво подписанные метрики то сразу отлуп
		encodeFunc := hash.CreateEncodeFunc(h.conf.SecretKey)
		for _, metric := range metrics {
			ok := hash.CheckSign(h.conf.SecretKey, metric.Hash, metric.CalcHash(encodeFunc))
			if !ok {
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write([]byte("incorrect metric sign"))
				return
			}
		}

		for _, metric := range metrics {
			err = h.repo.UpdateMetric(metric)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write([]byte("incorrect metric type"))
				return
			}
		}

		err = h.pStore.Save(r.Context())
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("failed to save metrics"))
			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write(AnswerOK)
	}
}

// ValueJSON обрабатывает POST запрос, который возвращает значение конкретной метрики в виде json
// @Tags Info
// @Summary Возвращает состояние метрики в формате json
// @Description На вход принимаем json с параметрами интересующей метрики в ответ шлем json с ее состоянием + подпись
// @ID infoValueJSON
// @Accept  json
// @Produce json
// @Param ID body string true "Metric name"
// @Param MType body string true "Metric type"
// @Success 200 {object} Metrics
// @Failure 404 {string} string "not found"
// @Failure 500 {string} string "content type not support, Read data problem, invalid deserialization"
// @Router /value/ [post]
func (h *Handler) ValueJSON() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("content type not support"))
			return
		}

		reqBody, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Read data problem"))
			return
		}

		metric := repository.Metrics{}
		err = json.Unmarshal(reqBody, &metric)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("invalid deserialization"))
			return
		}

		metricType := metric.MType
		metricName := metric.ID

		val, ok := h.repo.Metric(metricName, metricType)
		if ok {
			encodeFunc := hash.CreateEncodeFunc(h.conf.SecretKey)
			val.Hash = val.CalcHash(encodeFunc)

			result, ok := json.Marshal(val)
			if ok != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write([]byte("invalid serialization"))
				return
			}
			rw.Header().Set("Content-Type", "application/json")
			rw.Write(result)
			return
		}

		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("not found"))
	}
}

// Ping возвращает состояние доступности базы данных
// @Tags Info
// @Summary Запрос состояния доступности базы данных
// @ID infoPing
// @Produce plain
// @Success 200 {string} string "ok"
// @Failure 500 {string} string "ping error"
// @Router /ping [get]
func (h *Handler) Ping() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		err := h.pStore.Ping()
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("ping error"))
			return
		}
		rw.WriteHeader(http.StatusOK)
		rw.Write(AnswerOK)

	}
}
