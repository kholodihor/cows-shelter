package models

import (
	"time"

	"gorm.io/gorm"
)

type Partner struct {
    gorm.Model
    ID               uint      `json:"id" gorm:"primaryKey"`
    Name          	 string    `json:"name"`
    Logo          	 string    `json:"logo"`
	Link    		 string    `json:"link"`
     CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
}
