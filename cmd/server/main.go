// пакеты исполняемых приложений должны называться main
package main

import (
	"net/http"
	"strconv"
	"strings"
)

type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

// функция main вызывается автоматически при запуске приложения
func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

// функция run будет полезна при инициализации зависимостей сервера перед запуском
func run() error {
	store := MemStorage{}
	store.gauges = make(map[string]float64)
	store.counters = make(map[string]int64)

	return http.ListenAndServe(`:8080`, http.HandlerFunc(store.webhook))
}

// function to update data in store when they are valid
func (store *MemStorage) apply(metricType string, metricName string, metricValue string) error {
	if metricType == "counter" {
		i, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			return err
		}
		previous := store.counters[metricName]
		store.counters[metricName] = previous + i
		//fmt.Println("Added", i)
	} else {
		// gauge
		f, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return err
		}
		store.gauges[metricName] = f
		//fmt.Println("Replaced", f)
	}
	return nil
}

// функция webhook — обработчик HTTP-запроса
func (store *MemStorage) webhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// разрешаем только POST-запросы
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	line := r.RequestURI
	//fmt.Println(line)
	if !strings.HasPrefix(line, "/update/") {
		// must begin with update
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	split := strings.Split(line, "/")
	//fmt.Println(split[0], split[1])
	if split[2] != "gauge" && split[2] != "counter" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(split) != 5 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err := store.apply(split[2], split[3], split[4])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// установим правильный заголовок для типа данных
	w.Header().Set("Content-Type", "plain/text")
	// пока установим ответ-заглушку, без проверки ошибок
	_, _ = w.Write([]byte(`Applied`))
}
