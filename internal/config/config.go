package config

import (
    "flag"
    "os"
)

var FlagRunAddr string
var FlagBaseAddr string
var FlagLogLevel string
var FlagFileStoragePath string

func ParseFlags() {
    flag.StringVar(&FlagRunAddr, "a", "127.0.0.1:8080", "адрес запуска HTTP-сервера")
	flag.StringVar(&FlagBaseAddr, "b", "http://127.0.0.1:8080/qsd54gFg", "базовый адрес результирующего сокращённого URL")
	flag.StringVar(&FlagLogLevel, "l", "info", "log level")
	flag.StringVar(&FlagFileStoragePath, "f", "storage/storage.json", "путь до файла, куда сохраняются данные в формате JSON")
	
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
} 