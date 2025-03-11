package config

import (
	"flag"
	"os"
	"strconv"
)

// Флаги конфигурации для приложения, которые могут быть переданы через командную строку или переменные окружения.
// Эти флаги используются для настройки параметров, таких как адрес сервера, базовый URL, уровень логирования и другие.
var (
	// FlagRunAddr указывает адрес и порт, на котором должен запускаться HTTP-сервер.
	FlagRunAddr string

	// FlagBaseAddr указывает базовый адрес, который используется для формирования сокращенных URL.
	FlagBaseAddr string

	// FlagLogLevel указывает уровень логирования (например, "debug", "info", "warn", "error").
	FlagLogLevel string

	// FlagFileStoragePath указывает путь к файлу, в который сохраняются данные в формате JSON.
	FlagFileStoragePath string

	// FlagDatabaseDSN указывает строку подключения к базе данных.
	FlagDatabaseDSN string

	// FlagEnableHTTPS включает HTTPS в веб-сервере.
	FlagEnableHTTPS bool
)

// ParseFlags анализирует флаги командной строки и переменные окружения,
// чтобы установить значения для соответствующих переменных конфигурации.
func ParseFlags() {
	flag.StringVar(&FlagRunAddr, "a", "127.0.0.1:8080", "адрес запуска HTTP-сервера")
	flag.StringVar(&FlagBaseAddr, "b", "http://127.0.0.1:8080/qsd54gFg", "базовый адрес результирующего сокращённого URL")
	flag.StringVar(&FlagLogLevel, "l", "info", "log level")
	flag.StringVar(&FlagFileStoragePath, "f", "tmp/storage.json", "путь до файла, куда сохраняются данные в формате JSON")
	flag.StringVar(&FlagDatabaseDSN, "d", "", "адрес подключения к БД")
	flag.BoolVar(&FlagEnableHTTPS, "s", false, "включение HTTPS в веб-сервере")

	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		FlagRunAddr = envRunAddr
	}
	if envBaseAddr := os.Getenv("BASE_URL"); envBaseAddr != "" {
		FlagBaseAddr = envBaseAddr
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		FlagFileStoragePath = envFileStoragePath
	}
	if envDatabaseDSN := os.Getenv("DATABASE_DSN"); envDatabaseDSN != "" {
		FlagDatabaseDSN = envDatabaseDSN
	}
	if envEnableHTTPS := os.Getenv("ENABLE_HTTPS"); envEnableHTTPS != "" {
		if val, err := strconv.ParseBool(envEnableHTTPS); err == nil {
			FlagEnableHTTPS = val
		}
	}
}
