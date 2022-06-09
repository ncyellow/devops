package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ncyellow/devops/internal/server/storage"
)

// NewRouter создает chi.NewRouter и описывает маршрутизацию по url
func NewRouter(repo storage.Repository) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Get("/", ListHandler(repo))
	r.Get("/value/{metricType}/{metricName}", ValueHandler(repo))
	r.Post("/update/{metricType}/{metricName}/{metricValue}", UpdateHandler(repo))
	r.Post("/update/", UpdateJSONHandler(repo))
	r.Post("/update/", UpdateJSONHandler(repo))
	r.Post("/value/", ValueJSONHandler(repo))
	return r
}

// ValueHandler обрабатывает GET запросы на чтение значения метрик. Пример /value/counter/counterName
func ValueHandler(repo storage.Repository) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")
		switch metricType {
		case storage.Gauge:
			val, ok := repo.Gauge(metricName)
			if ok {
				rw.Write([]byte(fmt.Sprintf("%.03f", val)))
				return
			} else {
				rw.WriteHeader(http.StatusNotFound)
				return
			}
		case storage.Counter:
			val, ok := repo.Counter(metricName)
			if ok {
				rw.Write([]byte(fmt.Sprintf("%d", val)))
				return
			} else {
				rw.WriteHeader(http.StatusNotFound)
				return
			}
		default:
			rw.WriteHeader(http.StatusNotFound)
		}
	}
}

// UpdateHandler обрабатывает POST запросы на обновление значения метрик. Пример /update/counter/counterName/100
func UpdateHandler(repo storage.Repository) http.HandlerFunc {
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
		case storage.Gauge:
			value, err := strconv.ParseFloat(metricValue, 64)
			//! Второй параметр обязательно кастится в float64
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write([]byte("incorrect metric value"))
				return
			}
			err = repo.UpdateGauge(metricName, value)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write([]byte("incorrect metric name "))
				return
			}
		case storage.Counter:
			value, err := strconv.ParseInt(metricValue, 10, 64)
			//! Второй параметр обязательно кастится в int64
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write([]byte("incorrect metric value"))
				return
			}
			err = repo.UpdateCounter(metricName, value)
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
		rw.Write([]byte("ok"))
	}
}

// ListHandler обрабатывает GET запросы на корень url. Возвращает список всех метрик + значение
func ListHandler(repo storage.Repository) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(repo.String()))
	}
}

// UpdateJSONHandler обрабатывает POST запросы на обновление метрик в виде json
func UpdateJSONHandler(repo storage.Repository) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("content type not support"))
			return
		}
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Read data problem"))
			return
		}
		metric := storage.Metrics{}
		err = json.Unmarshal(reqBody, &metric)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("invalid deserialization"))
			return
		}

		err = repo.UpdateMetric(metric)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("incorrect metric type"))
			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("ok"))
	}
}

// ValueJSONHandler обрабатывает POST запрос, который возвращает список всех метрик в виде json
func ValueJSONHandler(repo storage.Repository) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Content-Type") != "application/json" {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("content type not support"))
			return
		}

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Read data problem"))
			return
		}

		metric := storage.Metrics{}
		err = json.Unmarshal(reqBody, &metric)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("invalid deserialization"))
			return
		}

		metricType := metric.MType
		metricName := metric.ID

		val, ok := repo.Metric(metricName, metricType)
		if ok {
			result, ok := json.Marshal(val)
			if ok != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write([]byte("invalid serialization"))
				return
			}
			rw.Write(result)
			return
		}
		rw.WriteHeader(http.StatusNotFound)
		return
	}
}
