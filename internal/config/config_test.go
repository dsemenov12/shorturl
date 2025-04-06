package config

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Тест загрузки конфигурации из JSON-файла
func TestParseFlags_ConfigFile(t *testing.T) {
	// Создаем временный конфигурационный файл
	configJSON := `{
		"server_address": "192.168.1.100:9090",
		"base_url": "http://192.168.1.100:9090/qsd54gFg",
		"file_storage_path": "/tmp/config_storage.json",
		"database_dsn": "postgres://user:password@localhost/db",
		"enable_https": true
	}`
	tempFile, err := os.CreateTemp("", "config.json")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write([]byte(configJSON))
	assert.NoError(t, err)
	tempFile.Close()

	// Устанавливаем путь к конфигурационному файлу
	os.Args = []string{"cmd", "-c", tempFile.Name()}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// Парсим флаги
	ParseFlags()

	// Проверяем, что параметры установлены из файла конфигурации
	assert.Equal(t, "http://127.0.0.1:8080/qsd54gFg", FlagBaseAddr)
	assert.Equal(t, "/tmp/config_storage.json", FlagFileStoragePath)
	assert.Equal(t, "postgres://user:password@localhost/db", FlagDatabaseDSN)
	assert.True(t, FlagEnableHTTPS)
}

// Тестируем ParseFlags с флагами командной строки
func TestParseFlags_CommandLineFlags(t *testing.T) {
	// Устанавливаем флаги командной строки
	os.Args = []string{"cmd", "-a", "127.0.0.1:9090", "-b", "http://localhost:9090/qsd54gFg", "-l", "debug", "-f", "/tmp/custom_storage.json", "-d", "user:password@tcp(localhost:3306)/db"}

	// Сброс флагов, чтобы избежать конфликта с флагами, определенными ранее
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// Парсим флаги
	ParseFlags()

	// Проверяем, что флаги установлены корректно
	assert.Equal(t, "127.0.0.1:9090", FlagRunAddr)
	assert.Equal(t, "http://localhost:9090/qsd54gFg", FlagBaseAddr)
	assert.Equal(t, "debug", FlagLogLevel)
	assert.Equal(t, "/tmp/custom_storage.json", FlagFileStoragePath)
	assert.Equal(t, "user:password@tcp(localhost:3306)/db", FlagDatabaseDSN)
}

// Тестируем ParseFlags с переменными окружения
func TestParseFlags_EnvironmentVariables(t *testing.T) {
	// Устанавливаем переменные окружения
	os.Setenv("SERVER_ADDRESS", "192.168.1.1:8080")
	os.Setenv("BASE_URL", "http://192.168.1.1:8080/qsd54gFg")
	os.Setenv("FILE_STORAGE_PATH", "/tmp/storage_from_env.json")
	os.Setenv("DATABASE_DSN", "postgres://user:password@localhost/db")

	// Устанавливаем значения флагов, чтобы они не влияли на тесты
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	FlagRunAddr = "default"
	FlagBaseAddr = "default"
	FlagLogLevel = "default"
	FlagFileStoragePath = "default"
	FlagDatabaseDSN = "default"

	// Парсим флаги
	ParseFlags()

	// Проверяем, что значения установлены из переменных окружения
	assert.Equal(t, "192.168.1.1:8080", FlagRunAddr)
	assert.Equal(t, "http://192.168.1.1:8080/qsd54gFg", FlagBaseAddr)
	assert.Equal(t, "/tmp/storage_from_env.json", FlagFileStoragePath)
	assert.Equal(t, "postgres://user:password@localhost/db", FlagDatabaseDSN)

	// Проверяем, что флаг логирования не был изменен переменной окружения
	//assert.Equal(t, "info", FlagLogLevel)

	// Очищаем переменные окружения после теста
	os.Unsetenv("SERVER_ADDRESS")
	os.Unsetenv("BASE_URL")
	os.Unsetenv("FILE_STORAGE_PATH")
	os.Unsetenv("DATABASE_DSN")
}

// Тестируем приоритет переменных окружения над флагами командной строки
func TestParseFlags_EnvironmentVariablesOverrideFlags(t *testing.T) {
	// Устанавливаем флаги командной строки
	os.Args = []string{"cmd", "-a", "127.0.0.1:9090", "-b", "http://localhost:9090/qsd54gFg"}

	// Сброс флагов, чтобы избежать конфликта с флагами, определенными ранее
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// Устанавливаем переменные окружения, которые должны переопределить флаги
	os.Setenv("SERVER_ADDRESS", "192.168.1.1:8080")
	os.Setenv("BASE_URL", "http://192.168.1.1:8080/qsd54gFg")

	// Парсим флаги
	ParseFlags()

	// Проверяем, что значения переменных окружения переопределили флаги
	assert.Equal(t, "192.168.1.1:8080", FlagRunAddr)
	assert.Equal(t, "http://192.168.1.1:8080/qsd54gFg", FlagBaseAddr)

	// Очищаем переменные окружения после теста
	os.Unsetenv("SERVER_ADDRESS")
	os.Unsetenv("BASE_URL")
}
