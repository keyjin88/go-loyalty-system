package config

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/keyjin88/go-loyalty-system/internal/app/logger"
	"log"
)

type Config struct {
	ServerAddress               string `env:"RUN_ADDRESS"`
	GinReleaseMode              bool   `env:"GIN_MODE"`
	LogLevel                    string `env:"LOG_LEVEL"`
	DataBaseURI                 string `env:"DATABASE_URI"`
	SecretKey                   string `env:"SECRET_KEY"`
	AccrualSystemAddress        string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	WorkerPoolSize              int    `env:"WORKER_POOL_SIZE"`
	ProcessingChannelBufferSize int    `env:"PROCESSING_CHANNEL_BUFFER_SIZE"`
}

func NewConfig() *Config {
	return &Config{}
}

// InitConfig обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func (config *Config) InitConfig() {
	flag.StringVar(&config.ServerAddress, "a", "localhost:8081", "address and port to run server")
	flag.BoolVar(&config.GinReleaseMode, "grm", false, "gin release mode")
	flag.StringVar(&config.LogLevel, "ll", "info", "log level")
	secretKey, genSecretErr := GenSecretKey(32)
	if genSecretErr != nil {
		secretKey = "abcdefghijklmnopqrstuvwxyz123456"
		logger.Log.Errorf("error while generate secret key")
	}
	flag.StringVar(&config.SecretKey, "sk", secretKey, "secret key for cryptographic")
	flag.StringVar(&config.AccrualSystemAddress, "r", "http://localhost:8080", "accrual system address")
	//flag.StringVar(&config.DataBaseURI, "d", "", "database dsn")
	// Оставил для локальных тестов
	flag.StringVar(&config.DataBaseURI, "d", "host=localhost user=pgadmin password=postgres dbname=loyaltydb port=5432 sslmode=disable", "database URI")
	flag.IntVar(&config.WorkerPoolSize, "wps", 10, "Worker pool size")
	flag.IntVar(&config.ProcessingChannelBufferSize, "pcbs", 10, "Processing channel buffer size")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
	// Пробуем распарсить переменные окружения, если их не будет, то оставляем значения по умолчанию из флагов
	err := env.Parse(config)
	if err != nil {
		log.Fatal(err)
	}
}

// GenSecretKey генерирует секретный ключ и возвращает его в виде Base64-строки.
func GenSecretKey(n int) (string, error) {
	data := make([]byte, n)
	_, err := rand.Read(data)
	if err != nil {
		return ``, err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}
