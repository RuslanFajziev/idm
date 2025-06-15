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

type RequestValidationError struct {
	Message string
}

type AlreadyExistsError struct {
	Message string
}

func (err RequestValidationError) Error() string {
	return err.Message
}

func (err AlreadyExistsError) Error() string {
	return err.Message
}

// GetConfig получение конфигурации из .env файла или переменных окружения
func GetConfig(envFile string) (Config, error) {
	err := godotenv.Load(envFile)

	if err != nil {
		return Config{}, fmt.Errorf("warning: Could not load .env file: %v\n", err)
	}

	cfg := Config{
		DbDriverName: os.Getenv("DB_DRIVER_NAME"),
		Dsn:          os.Getenv("DB_DSN"),
	}

	if cfg.DbDriverName == "" || cfg.Dsn == "" {
		return Config{}, fmt.Errorf("DB_DRIVER_NAME and DB_DSN must be set in .env file or environment variables")
	}

	return cfg, nil
}
