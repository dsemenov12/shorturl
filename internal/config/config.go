package config

import (
    "flag"
)

var FlagRunAddr string
var FlagBaseAddr string

func ParseFlags() {
    flag.StringVar(&FlagRunAddr, "a", "127.0.0.1:8080", "адрес запуска HTTP-сервера")
	flag.StringVar(&FlagBaseAddr, "b", "http://localhost:8000/qsd54gFg", "базовый адрес результирующего сокращённого URL")
	
    flag.Parse()
} 