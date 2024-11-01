package config

import (
    "flag"
    "os"
)

var FlagRunAddr string
var FlagBaseAddr string

func ParseFlags() {
    flag.StringVar(&FlagRunAddr, "a", "127.0.0.1:8080", "адрес запуска HTTP-сервера")
	flag.StringVar(&FlagBaseAddr, "b", "http://127.0.0.1:8080/qsd54gFg", "базовый адрес результирующего сокращённого URL")
	
    flag.Parse()

    if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
        FlagRunAddr = envRunAddr
    }
    if envBaseAddr := os.Getenv("BASE_URL"); envBaseAddr != "" {
        FlagBaseAddr = envBaseAddr
    }
} 