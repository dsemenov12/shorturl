package logger

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestBuildLogger(t *testing.T) {
	logger, err := buildLogger("info")
	assert.NoError(t, err)
	assert.NotNil(t, logger)
}

func TestRequestLogger(t *testing.T) {
	// Мокаем логгер для проверки выводимых логов
	logger := zap.NewNop()

	// Устанавливаем глобальный логгер
	Log = logger

	// Создаем простого обработчика, который просто отвечает на запрос
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}

	// Создаем тестовый запрос
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	// Создаем логирующий обработчик
	handlerWithLogging := RequestLogger(handler)

	// Запускаем обработчик
	handlerWithLogging(rr, req)

	// Проверяем результат
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "OK", rr.Body.String())
}

func TestLoggingResponseWriter(t *testing.T) {
	// Создаем responseWriter с логированием
	responseData := &responseData{}
	lrw := loggingResponseWriter{
		ResponseWriter: httptest.NewRecorder(),
		responseData:   responseData,
	}

	// Проверяем работу метода Write
	n, err := lrw.Write([]byte("Test Response"))
	assert.NoError(t, err)
	assert.Equal(t, 13, n)
	assert.Equal(t, 13, responseData.size)

	// Проверяем работу метода WriteHeader
	lrw.WriteHeader(http.StatusOK)
	assert.Equal(t, http.StatusOK, responseData.status)
}

func TestRequestLoggerLogsCorrectly(t *testing.T) {
	// Мокаем логгер для проверки выводимых логов
	logger := zap.NewNop()

	// Устанавливаем глобальный логгер
	Log = logger

	// Создаем обработчик, который просто отвечает на запрос
	handler := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}

	// Создаем тестовый запрос
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	// Создаем логирующий обработчик
	handlerWithLogging := RequestLogger(handler)

	// Запускаем обработчик
	handlerWithLogging(rr, req)

	// Проверяем результат
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "OK", rr.Body.String())
}
