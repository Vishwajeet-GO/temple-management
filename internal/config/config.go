package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort    string
	AppEnv     string
	DBType     string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	return &Config{
		AppPort:    getEnv("APP_PORT", "8080"),
		AppEnv:     getEnv("APP_ENV", "development"),
		DBType:     getEnv("DB_TYPE", "sqlite"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "temple.db"),
		JWTSecret:  getEnv("JWT_SECRET", "default-secret-key"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
