package config

import (
    "flag"
    "os"
)

var FlagRunAddr string
var FlagBaseAddr string
var FlagLogLevel string
var FlagFileStoragePath string
var FlagDatabaseDSN string

func ParseFlags() {
    flag.StringVar(&FlagRunAddr, "a", "127.0.0.1:8080", "адрес запуска HTTP-сервера")
	flag.StringVar(&FlagBaseAddr, "b", "http://127.0.0.1:8080/qsd54gFg", "базовый адрес результирующего сокращённого URL")
	flag.StringVar(&FlagLogLevel, "l", "info", "log level")
	flag.StringVar(&FlagFileStoragePath, "f", "storage/storage.json", "путь до файла, куда сохраняются данные в формате JSON")
	flag.StringVar(&FlagDatabaseDSN, "d", "host=go-shortner-pg user=user password=123456 dbname=pg-go-shortner sslmode=disable", "адрес подключения к БД")
	
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
} 