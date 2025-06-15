package models

import (
	"time"

	"gorm.io/gorm"
)

type Review struct {
    gorm.Model
    ID               uint      `json:"id" gorm:"primaryKey"`
    NameEn           string    `json:"name_en"`
    NameUa           string    `json:"name_ua"`
    ReviewEn    	 string    `json:"review_en"`
    ReviewUa    	 string    `json:"review_ua"`
    CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
}
