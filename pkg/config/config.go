package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// DBConfig holds database connection settings
type DBConfig struct {
	Host     string `json:"host" env:"DB_HOST"`
	User     string `json:"user" env:"DB_USER"`
	Password string `json:"password" env:"DB_PASSWORD"`
	DBName   string `json:"dbname" env:"DB_NAME"`
	Port     string `json:"port" env:"DB_PORT"`
	SSLMode  string `json:"sslmode" env:"DB_SSLMODE"`
}

// LoadDBConfig loads the database config from environment variables
func LoadDBConfig() (*DBConfig, error) {
	_ = godotenv.Load() // Load .env file if present

	cfg := &DBConfig{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	// Set defaults if not provided
	if cfg.Host == "" {
		cfg.Host = "localhost"
	}
	if cfg.Port == "" {
		cfg.Port = "5432"
	}
	if cfg.SSLMode == "" {
		cfg.SSLMode = "disable"
	}

	// Validation (optional)
	if cfg.User == "" || cfg.Password == "" || cfg.DBName == "" {
		return nil, fmt.Errorf("missing required database config fields")
	}

	return cfg, nil
}
