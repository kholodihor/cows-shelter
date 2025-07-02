package models

import "gorm.io/gorm"

type Gallery struct {
    gorm.Model
    ImageUrl string `json:"image_url"`
}