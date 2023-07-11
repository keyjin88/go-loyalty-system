package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

type Config struct {
	ServerAddress        string `env:"RUN_ADDRESS"`
	GinReleaseMode       bool   `env:"GIN_MODE"`
	LogLevel             string `env:"LOG_LEVEL"`
	DataBaseURI          string `env:"DATABASE_URI"`
	SecretKey            string `env:"SECRET_KEY"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
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
	flag.StringVar(&config.SecretKey, "sk", "abcdefghijklmnopqrstuvwxyz123456", "secret key for cryptographic")
	flag.StringVar(&config.AccrualSystemAddress, "r", "http://localhost:8080", "accrual system address")
	//flag.StringVar(&config.DataBaseURI, "d", "", "database dsn")
	// Оставил для локальных тестов
	flag.StringVar(&config.DataBaseURI, "d", "host=localhost user=pgadmin password=postgres dbname=loyaltydb port=5432 sslmode=disable", "database URI")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
	// Пробуем распарсить переменные окружения, если их не будет, то оставляем значения по умолчанию из флагов
	err := env.Parse(config)
	if err != nil {
		log.Fatal(err)
	}
}
