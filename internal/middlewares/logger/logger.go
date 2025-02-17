package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Log является глобальной переменной для логгера, инициализированного с помощью пакета zap.
var Log *zap.Logger = zap.NewNop()

type (
	// Структура responseData хранит информацию о статусе и размере HTTP-ответа.
	responseData struct {
		status int // Код статуса ответа
		size   int // Размер ответа в байтах
	}

	// Структура loggingResponseWriter расширяет http.ResponseWriter и позволяет логировать данные о статусе и размере ответа.
	loggingResponseWriter struct {
		http.ResponseWriter               // Встроенный ResponseWriter для передачи данных в стандартный поток
		responseData        *responseData // Ссылка на структуру, содержащую статус и размер ответа
	}
)

// Write перезаписывает метод Write для записи в ответ и отслеживания размера ответа.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// WriteHeader перезаписывает метод WriteHeader для записи статуса HTTP-ответа.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// Initialize инициализирует глобальный логгер Log с заданным уровнем логирования.
// Принимает строковое значение уровня логирования (например, "debug", "info", "error").
// Возвращает ошибку, если уровень не может быть разобран или если возникли проблемы при создании логгера.
func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	zl, err := cfg.Build()
	if err != nil {
		return err
	}

	Log = zl
	return nil
}

// RequestLogger является middleware для логирования информации о входящих HTTP-запросах и их ответах.
// Логируются метод, путь запроса, продолжительность запроса, код статуса ответа и размер ответа.
// handlerFunc: обработчик HTTP-запроса, который будет обернут в логирование.
//
// Возвращаемое значение: возвращает обработчик HTTP-запроса, который логирует информацию о запросах и ответах.
func RequestLogger(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		handlerFunc(&lw, r)

		duration := time.Since(start)

		Log.Info("got incoming HTTP request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Duration("duration", duration),
		)
		Log.Info("got response",
			zap.Int("status", responseData.status),
			zap.Int("size", responseData.size),
		)
	})
}
