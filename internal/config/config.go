package config

import (
	"encoding/json"
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

	// FlagConfigFilePath путь к файлу конфигурации.
	FlagConfigFilePath string

	// FlagTrustedSubnet указывает доверенную подсеть в формате CIDR.
	FlagTrustedSubnet string

	FlagGRPCAddress string

	FlagGRPCGatewayAddr string

	FlagEnableGRPCGateway bool
)

// Config структура для JSON-конфигурации
type Config struct {
	ServerAddress      string `json:"server_address"`
	BaseURL            string `json:"base_url"`
	FileStoragePath    string `json:"file_storage_path"`
	DatabaseDSN        string `json:"database_dsn"`
	EnableHTTPS        bool   `json:"enable_https"`
	TrustedSubnet      string `json:"trusted_subnet"`
	GRPCAddress        string `json:"grpc_address"`
	GRPCGatewayAddress string `json:"grpc_gateway_address"`
	EnableGRPCGateway  bool   `json:"enable_grpc_gateway"`
}

// ParseFlags анализирует флаги командной строки и переменные окружения,
// чтобы установить значения для соответствующих переменных конфигурации.
func ParseFlags() {
	flag.StringVar(&FlagRunAddr, "a", "127.0.0.1:8080", "адрес запуска HTTP-сервера")
	flag.StringVar(&FlagBaseAddr, "b", "http://127.0.0.1:8080/qsd54gFg", "базовый адрес результирующего сокращённого URL")
	flag.StringVar(&FlagLogLevel, "l", "info", "log level")
	flag.StringVar(&FlagFileStoragePath, "f", "tmp/storage.json", "путь до файла, куда сохраняются данные в формате JSON")
	flag.StringVar(&FlagDatabaseDSN, "d", "", "адрес подключения к БД")
	flag.BoolVar(&FlagEnableHTTPS, "s", false, "включение HTTPS в веб-сервере")
	flag.StringVar(&FlagConfigFilePath, "c", "", "путь до JSON-файла конфигурации")
	flag.StringVar(&FlagConfigFilePath, "config", "", "путь до JSON-файла конфигурации (аналог -c)")
	flag.StringVar(&FlagTrustedSubnet, "t", "", "доверенная подсеть в формате CIDR")
	flag.StringVar(&FlagGRPCAddress, "grpc-address", "127.0.0.1:9090", "адрес запуска gRPC-сервера")
	flag.StringVar(&FlagGRPCGatewayAddr, "grpc-gateway-address", "127.0.0.1:8081", "адрес запуска grpc-gateway HTTP сервера")
	flag.BoolVar(&FlagEnableGRPCGateway, "enable-grpc-gateway", false, "включить HTTP/REST gRPC-Gateway")

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
	if envConfigFilePath := os.Getenv("CONFIG"); envConfigFilePath != "" {
		FlagConfigFilePath = envConfigFilePath
	}
	if envTrustedSubnet := os.Getenv("TRUSTED_SUBNET"); envTrustedSubnet != "" {
		FlagTrustedSubnet = envTrustedSubnet
	}
	if envGRPCAddress := os.Getenv("GRPC_ADDRESS"); envGRPCAddress != "" {
		FlagGRPCAddress = envGRPCAddress
	}
	if envGRPCGatewayAddr := os.Getenv("GRPC_GATEWAY_ADDRESS"); envGRPCGatewayAddr != "" {
		FlagGRPCGatewayAddr = envGRPCGatewayAddr
	}
	if envEnableGRPCGateway := os.Getenv("ENABLE_GRPC_GATEWAY"); envEnableGRPCGateway != "" {
		if val, err := strconv.ParseBool(envEnableGRPCGateway); err == nil {
			FlagEnableGRPCGateway = val
		}
	}

	if FlagConfigFilePath != "" {
		loadConfigFromFile(FlagConfigFilePath)
	}
}

// loadConfigFromFile загружает конфигурацию из JSON-файла
func loadConfigFromFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		return // Если файл не найден, просто продолжаем с текущими настройками
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return
	}

	// Устанавливаем значения из файла, только если они не переопределены
	if FlagRunAddr == "127.0.0.1:8080" {
		FlagRunAddr = cfg.ServerAddress
	}
	if FlagBaseAddr == "http://127.0.0.1:8080" {
		FlagBaseAddr = cfg.BaseURL
	}
	if FlagFileStoragePath == "tmp/storage.json" {
		FlagFileStoragePath = cfg.FileStoragePath
	}
	if FlagDatabaseDSN == "" {
		FlagDatabaseDSN = cfg.DatabaseDSN
	}
	if !FlagEnableHTTPS { // Если по умолчанию false, заменяем значением из конфига
		FlagEnableHTTPS = cfg.EnableHTTPS
	}
	if FlagTrustedSubnet == "" {
		FlagTrustedSubnet = cfg.TrustedSubnet
	}
	if FlagGRPCAddress == "127.0.0.1:9090" {
		FlagGRPCAddress = cfg.GRPCAddress
	}
	if FlagGRPCGatewayAddr == "127.0.0.1:8081" {
		FlagGRPCGatewayAddr = cfg.GRPCGatewayAddress
	}
	if !FlagEnableGRPCGateway {
		FlagEnableGRPCGateway = cfg.EnableGRPCGateway
	}
}
