package models

import (
	"time"

	"gorm.io/gorm"
)

type Pdf struct {
    gorm.Model
    ID               	uint      `json:"id" gorm:"primaryKey"`
    Title         		string    `json:"title"`
    DocumentUrl         string    `json:"document_url"`
    CreatedAt           time.Time `json:"created_at" gorm:"autoCreateTime"`
}