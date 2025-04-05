package models

import (
	"time"
)

// User represents a user entity in the database
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Fullname  string    `gorm:"type:varchar(255);not null" json:"full_name"`
	Username  *string   `gorm:"type:varchar(100);unique;not null" json:"username"`
	Password  string    `gorm:"type:varchar(255);not null" json:"password"` // Hashed password
	CreatedAt time.Time `json:"created_at"`
}
