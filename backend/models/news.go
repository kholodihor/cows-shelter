package models

import (
	"time"

	"gorm.io/gorm"
)

type News struct {
    gorm.Model
    ID               uint      `json:"id" gorm:"primaryKey"`
    TitleEn          string    `json:"title_en"`
    TitleUa          string    `json:"title_ua"`
    SubtitleEn    	 string    `json:"subtitle_en"`
    SubtitleUa       string    `json:"subtitle_ua"`
    ContentEn    	 string    `json:"content_en"`
    ContentUa        string     `json:"content_ua"`
    ImageUrl         string    `json:"image_url"`
    CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
}
