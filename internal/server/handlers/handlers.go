package handlers

import (
	"github.com/ncyellow/devops/internal/server/storage"
	"net/http"
	"strconv"
	"strings"
)

func Handler(repo storage.Repository) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		args := strings.Split(strings.TrimPrefix(r.URL.Path, "/update/"), "/")

		//! Метод только post
		if r.Method != http.MethodPost {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("only post method support"))
			return
		}

		//! Мы должны иметь три параметра все остальное отлуп
		if len(args) != 3 {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte("incorrect format"))
			return
		}

		metricType := args[0]
		name := args[1]

		switch metricType {
		case "gauge":
			value, err := strconv.ParseFloat(args[2], 64)
			//! Второй параметр обязательно кастится в float64
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write([]byte("incorrect metric value"))
				return
			}
			err = repo.UpdateGauge(name, value)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write([]byte("incorrect metric name "))
				return
			}
		case "counter":
			value, err := strconv.ParseInt(args[2], 10, 64)
			//! Второй параметр обязательно кастится в int64
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write([]byte("incorrect metric value"))
				return
			}
			err = repo.UpdateCounter(name, value)
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
