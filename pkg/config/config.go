package config

// DBConfig holds database connection settings
type DBConfig struct {
	Host     string `json:"host" env:"DB_HOST"`
	User     string `json:"user" env:"DB_USER"`
	Password string `json:"password" env:"DB_PASSWORD"`
	DBName   string `json:"dbname" env:"DB_NAME"`
	Port     string `json:"port" env:"DB_PORT"`
	SSLMode  string `json:"sslmode" env:"DB_SSLMODE"`
}
