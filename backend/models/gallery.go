package models

import (
	"time"

	"gorm.io/gorm"
)

type Gallery struct {
    gorm.Model
    ID               uint      `json:"id" gorm:"primaryKey"`
    ImageUrl         string    `json:"image_url"`
    CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
}