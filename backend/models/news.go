package models

import "gorm.io/gorm"

type News struct {
    gorm.Model
    TitleEn     string `json:"title_en"`
    TitleUa     string `json:"title_ua"`
    SubtitleEn  string `json:"subtitle_en"`
    SubtitleUa  string `json:"subtitle_ua"`
    ContentEn   string `json:"content_en"`
    ContentUa   string `json:"content_ua"`
    ImageUrl    string `json:"image_url"`
}
