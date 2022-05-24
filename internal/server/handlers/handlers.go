package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ncyellow/devops/internal/server/storage"
	"net/http"
	"strconv"
)

func NewRouter(repo storage.Repository) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Get("/", MetricListHandler(repo))
	r.Get("/value/{metricType}/{metricName}", MetricValueHandler(repo))
	r.Post("/update/{metricType}/{metricName}/{metricValue}", MetricUpdateHandler(repo))
	return r
}

func MetricValueHandler(repo storage.Repository) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")
		switch metricType {
		case storage.Gauge:
			val, ok := repo.Gauge(metricName)
			if ok {
				rw.Write([]byte(fmt.Sprintf("%f", val)))
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

func MetricUpdateHandler(repo storage.Repository) http.HandlerFunc {
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
				rw.Write([]byte("incorrect metric name "))
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

func MetricListHandler(repo storage.Repository) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("this is all metrics"))
	}
}
