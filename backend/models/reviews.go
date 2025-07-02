package models

import "gorm.io/gorm"

type Review struct {
    gorm.Model
    NameEn   string `json:"name_en"`
    NameUa   string `json:"name_ua"`
    ReviewEn string `json:"review_en"`
    ReviewUa string `json:"review_ua"`
}
