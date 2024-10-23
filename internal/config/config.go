package config

import (
    "flag"
)

var FlagRunAddr string
var FlagBaseAddr string

func ParseFlags() {
    flag.StringVar(&FlagRunAddr, "a", ":8080", "адрес запуска HTTP-сервера")
	flag.StringVar(&FlagBaseAddr, "b", "", "базовый адрес результирующего сокращённого URL")
	
    flag.Parse()
} 