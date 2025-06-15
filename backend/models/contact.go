package models

import (
	"time"
)

type Contact struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"not null;uniqueIndex" json:"email"`
	Phone     string    `gorm:"not null" json:"phone"`
	ImageURL  string    `gorm:"type:text" json:"imageUrl,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
