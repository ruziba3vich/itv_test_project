package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		DBConfig   *DBConfig
		Redis      *RedisConfig
		JwtSecret  string
		RLConfig   *RateLimiterConfig
		AppPort    string
		AccessTTL  int
		RefreshTTL int
		MovieTTL   int
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

	RateLimiterConfig struct {
		MaxTokens  int
		RefillRate float64
		Window     time.Duration
	}
)

// DBConfig holds database connection settings

// LoadDBConfig loads the database config from environment variables
func LoadConfig() *Config {
	_ = godotenv.Load() // Load .env file if present

	cfg := &Config{
		DBConfig: &DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			User:     getEnv("DB_USER", "root_user"),
			Password: getEnv("DB_PASSWORD", "Dost0n1k"),
			DBName:   getEnv("DB_NAME", "itv_test"),
			Port:     getEnv("DB_PORT", "5432"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Redis: &RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PWD", "password"),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		JwtSecret: getEnv("JWT_SECRET", "prodonik"),
		RLConfig: &RateLimiterConfig{
			MaxTokens:  getEnvInt("RL_MAX_TOKENS", 4),
			Window:     time.Duration(getEnvInt("RL_WINDOW", 1) * int(time.Minute)),
			RefillRate: getEnvFloat("RL_REFILL_RATE", 0.25),
		},
		AppPort:    getEnv("APP_PORT", ":7777"),
		AccessTTL:  getEnvInt("ACCESS_TTL", 15),
		RefreshTTL: getEnvInt("REFRESH_TTL", 30),
		MovieTTL:   getEnvInt("MOVIE_TTL", 20),
	}
	return cfg
}

func getEnvFloat(key string, fallback float64) float64 {
	if value, exists := os.LookupEnv(key); exists {
		window, _ := strconv.ParseFloat(value, 32)
		return window
	}
	return fallback
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
