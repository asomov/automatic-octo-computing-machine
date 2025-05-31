package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemStorage_webhook(t *testing.T) {
	tests := []struct {
		name  string
		input string
		code  int
	}{
		//positive tests
		{
			name:  "counter positive test",
			input: "/update/counter/metric1/1",
			code:  200,
		},
		{
			name:  "gauge positive test",
			input: "/update/gauge/metric1/1",
			code:  200,
		},
		// negative tests
		// counter
		{
			name:  "invalid type test",
			input: "/update/invalid/metric1/1",
			code:  400,
		},
		{
			name:  "missing metric name test",
			input: "/update/counter/1",
			code:  404,
		},
		{
			name:  "missing metric name test",
			input: "/update/counter/name2/none",
			code:  400,
		},
		{
			name:  "missing metric name test",
			input: "/update/counter/name2/10.0",
			code:  400,
		},
		// negative tests
		// gauge
		{
			name:  "invalid type test",
			input: "/update/invalid/metric1/1",
			code:  400,
		},
		{
			name:  "missing metric name test",
			input: "/update/gauge/1",
			code:  404,
		},
		{
			name:  "missing metric name test",
			input: "/update/gauge/name2/none",
			code:  400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &MemStorage{}
			store.gauges = make(map[string]float64)
			store.counters = make(map[string]int64)
			request := httptest.NewRequest(http.MethodPost, tt.input, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()

			store.webhook(w, request)
			res := w.Result()
			// проверяем код ответа
			require.Equal(t, tt.code, res.StatusCode, tt.name)
		})
	}
}

func Test_GET_webhook(t *testing.T) {
	tests := []struct {
		name  string
		input string
		code  int
	}{
		{
			name:  "first test",
			input: "/update/counter/metric1/1",
			code:  405,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &MemStorage{}
			store.gauges = make(map[string]float64)
			store.counters = make(map[string]int64)
			request := httptest.NewRequest(http.MethodGet, tt.input, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()

			store.webhook(w, request)
			res := w.Result()
			// проверяем код ответа
			require.Equal(t, tt.code, res.StatusCode)
		})
	}
}
