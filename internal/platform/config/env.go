package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	AppPort         string
	DatabaseURL     string
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	DBSSLMode       string
	JWTSecret       string
	JWTExpiresHours int
	AdminUsername   string
	AdminPassword   string
	AdminRole       string
}

func Load() (*EnvConfig, error) {
	if err := godotenv.Overload(); err != nil {
		return nil, fmt.Errorf("no se pudo cargar .env: %w", err)
	}

	expiresHoursStr := getEnvOrDefault("JWT_EXPIRES_HOURS", "24")
	expiresHours, err := strconv.Atoi(expiresHoursStr)
	if err != nil {
		return nil, fmt.Errorf("JWT_EXPIRES_HOURS invalido")
	}

	cfg := &EnvConfig{
		AppPort:         getEnvOrDefault("APP_PORT", "8080"),
		DatabaseURL:     getEnvOrDefault("DATABASE_URL", ""),
		DBHost:          getEnvOrDefault("DB_HOST", ""),
		DBPort:          getEnvOrDefault("DB_PORT", ""),
		DBUser:          getEnvOrDefault("DB_USER", ""),
		DBPassword:      getEnvOrDefault("DB_PASSWORD", ""),
		DBName:          getEnvOrDefault("DB_NAME", ""),
		DBSSLMode:       getEnvOrDefault("DB_SSLMODE", ""),
		JWTSecret:       getRequiredEnv("JWT_SECRET"),
		JWTExpiresHours: expiresHours,
		AdminUsername:   getRequiredEnv("ADMIN_USERNAME"),
		AdminPassword:   getRequiredEnv("ADMIN_PASSWORD"),
		AdminRole:       getEnvOrDefault("ADMIN_ROLE", "admin"),
	}

	if cfg.JWTSecret == "" || cfg.AdminUsername == "" || cfg.AdminPassword == "" {
		return nil, fmt.Errorf("faltan variables obligatorias de autenticacion en el archivo .env")
	}

	// Si no existe DATABASE_URL, exigimos las variables DB_* normales
	if cfg.DatabaseURL == "" {
		if cfg.DBHost == "" || cfg.DBPort == "" || cfg.DBUser == "" ||
			cfg.DBPassword == "" || cfg.DBName == "" || cfg.DBSSLMode == "" {
			return nil, fmt.Errorf("faltan variables obligatorias de base de datos")
		}
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

func getEnvOrDefault(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return defaultValue
	}
	return value
}
