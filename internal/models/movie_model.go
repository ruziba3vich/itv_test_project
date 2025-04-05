package models

import (
	"time"

	"gorm.io/gorm"
)

// Movie represents the movie entity in the database
type Movie struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"type:varchar(255);not null" json:"title"`
	Director  string         `gorm:"type:varchar(100);not null" json:"director"`
	Year      int            `gorm:"not null" json:"year"`
	Plot      string         `gorm:"type:text" json:"plot"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"` // Soft delete support
}
