package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type DBConfig struct {
	Host     string `envconfig:"HOST" required:"true"`
	Port     string `envconfig:"PORT" default:"5432"`
	User     string `envconfig:"USER" required:"true"`
	Password string `envconfig:"PASSWORD" required:"true"`
	Database string `envconfig:"DB" required:"true"`
}

func NewDBConfig() (DBConfig, error) {
	var config DBConfig

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found. Assuming environment variables are set (e.g., in Docker).")
	}

	if err := envconfig.Process("POSTGRES", &config); err != nil {
		return DBConfig{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}
