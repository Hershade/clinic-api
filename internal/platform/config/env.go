package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	AppPort    string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func Load() (*EnvConfig, error) {
	//intentar cargar el .env
	_ = godotenv.Load()

	cfg := &EnvConfig{
		AppPort:    getEnvOrDefault("APP_PORT", "8080"),
		DBHost:     getRequiredEnv("DB_HOST"),
		DBPort:     getRequiredEnv("DB_PORT"),
		DBUser:     getRequiredEnv("DB_USER"),
		DBPassword: getRequiredEnv("DB_PASSWORD"),
		DBName:     getRequiredEnv("DB_NAME"),
		DBSSLMode:  getRequiredEnv("DB_SSLMODE"),
	}

	if cfg.DBHost == "" || cfg.DBPort == "" || cfg.DBUser == "" ||
		cfg.DBPassword == "" || cfg.DBName == "" || cfg.DBSSLMode == "" {
		return nil, fmt.Errorf("Faltan variables obligatorias en el archivo .env")
	}

	return cfg, nil
}

func getRequiredEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return ""
	}
	return value
}

func getEnvOrDefault(key, defaulValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return defaulValue
	}
	return value
}
