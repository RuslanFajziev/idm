package common

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config общая конфигурация всего приложения
type Config struct {
	DbDriverName string `validate:"required"`
	Dsn          string `validate:"required"`
}

// GetConfig получение конфигурации из .env файла или переменных окружения
func GetConfig(envFile string) Config {
	err := godotenv.Load(envFile)

	if err != nil {
		fmt.Printf("Warning: Could not load .env file: %v", err)
	}

	cfg := Config{
		DbDriverName: os.Getenv("DB_DRIVER_NAME"),
		Dsn:          os.Getenv("DB_DSN"),
	}

	if cfg.DbDriverName == "" || cfg.Dsn == "" {
		panic("DB_DRIVER_NAME and DB_DSN must be set in .env file or environment variables")
	}

	return cfg
}
