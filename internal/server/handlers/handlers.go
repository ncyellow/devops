package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v4"
	"github.com/ncyellow/devops/internal/hash"
	"github.com/ncyellow/devops/internal/server/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ncyellow/devops/internal/server/middlewares"
	"github.com/ncyellow/devops/internal/server/storage"
)

type Handler struct {
	*chi.Mux
	conf   *config.Config
	repo   storage.Repository
	pStore storage.PersistentStorage
}

// NewRouter создает chi.NewRouter и описывает маршрутизацию по url
func NewRouter(repo storage.Repository, conf *config.Config, pStore storage.PersistentStorage) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middlewares.EncoderGZIP)

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
	r.Post("/ping", handler.Ping())

	return handler
}

func (h *Handler) List() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "text/html")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(h.repo.String()))
	}
}

func (h *Handler) Value() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")
		switch metricType {
		case storage.Gauge:
			val, ok := h.repo.Gauge(metricName)
			if ok {
				rw.Write([]byte(fmt.Sprintf("%.03f", val)))
				return
			} else {
				rw.WriteHeader(http.StatusNotFound)
				return
			}
		case storage.Counter:
			val, ok := h.repo.Counter(metricName)
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
		case storage.Gauge:
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
		case storage.Counter:
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
		rw.Write([]byte("ok"))
	}
}

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
		metric := storage.Metrics{}
		err = json.Unmarshal(reqBody, &metric)

		fmt.Printf("пришло - %#v\n", metric)
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

		fmt.Printf("Обновление метрики")
		err = h.repo.UpdateMetric(metric)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("incorrect metric type"))
			return
		}
		h.pStore.Save()

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("ok"))
	}
}

// ValueJSON обрабатывает POST запрос, который возвращает список всех метрик в виде json
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

		metric := storage.Metrics{}
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
	}
}

// Ping возвращает доступность базы данных
func (h *Handler) Ping() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		conn, err := pgx.Connect(context.Background(), h.conf.DatabaseConn)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			rw.Write([]byte("no connection"))
			return
		}
		defer conn.Close(context.Background())

		var result int
		err = conn.QueryRow(r.Context(), "select 1").Scan(&result)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("query error"))
			return
		}
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("ok"))
	}
}
