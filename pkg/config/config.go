package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		DBConfig  *DBConfig
		Redis     *RedisConfig
		JwtSecret string
	}

	RedisConfig struct {
		Host, Port, Password string
		DB                   int
	}

	DBConfig struct {
		Host     string `json:"host" env:"DB_HOST"`
		User     string `json:"user" env:"DB_USER"`
		Password string `json:"password" env:"DB_PASSWORD"`
		DBName   string `json:"dbname" env:"DB_NAME"`
		Port     string `json:"port" env:"DB_PORT"`
		SSLMode  string `json:"sslmode" env:"DB_SSLMODE"`
	}
)

// DBConfig holds database connection settings

// LoadDBConfig loads the database config from environment variables
func LoadConfig() *Config {
	_ = godotenv.Load() // Load .env file if present

	cfg := &Config{
		DBConfig: &DBConfig{
			Host:     getEnv(os.Getenv("DB_HOST"), "localhost"),
			User:     getEnv(os.Getenv("DB_USER"), "5432"),
			Password: getEnv(os.Getenv("DB_PASSWORD"), "password"),
			DBName:   getEnv(os.Getenv("DB_NAME"), "itv_test"),
			Port:     getEnv(os.Getenv("DB_PORT"), "5432"),
			SSLMode:  getEnv(os.Getenv("DB_SSLMODE"), "disable"),
		},
		Redis: &RedisConfig{
			Host:     getEnv(os.Getenv("REDIS_HOST"), "localhost"),
			Port:     getEnv(os.Getenv("REDIS_PORT"), "6379"),
			Password: getEnv(os.Getenv("REDIS_PWD"), "password"),
			DB:       getEnvInt(os.Getenv("REDIS_DB"), 0),
		},
		JwtSecret: getEnv("JWT_SECRET", "prodonik"),
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		var intValue int
		_, err := fmt.Sscanf(value, "%d", &intValue)
		if err == nil {
			return intValue
		}
	}
	return fallback
}
