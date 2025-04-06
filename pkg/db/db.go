package db

import (
	"fmt"

	"github.com/ruziba3vich/itv_test_project/internal/models"
	"github.com/ruziba3vich/itv_test_project/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDB initializes a new GORM database connection using the provided config
func NewDB(cfg *config.Config) (*gorm.DB, error) {
	// Construct the DSN from the config
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBConfig.Host,
		cfg.DBConfig.User,
		cfg.DBConfig.Password,
		cfg.DBConfig.DBName,
		cfg.DBConfig.Port,
		cfg.DBConfig.SSLMode,
	)

	// Open the database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Auto-migrate the schema
	if err := db.AutoMigrate(&models.Movie{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	if err := db.AutoMigrate(&models.RefreshToken{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}
	return db, nil
}
