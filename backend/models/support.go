package models

import (
	"time"

	"gorm.io/gorm"
)

type SupportCard struct {
    gorm.Model
    ID        uint      `json:"id" gorm:"primaryKey"`
    Title     string    `json:"title"`
    Subtitle  string    `json:"subtitle"`
    Banner    string    `json:"banner"`
    Image     string    `json:"image"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type SupportStep struct {
    gorm.Model
    ID        uint      `json:"id" gorm:"primaryKey"`
    StepText  string    `json:"step_text"`
    Order     int       `json:"order"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
