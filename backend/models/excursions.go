package models

import (
	"time"

	"gorm.io/gorm"
)

type Excursion struct {
    gorm.Model
    ID               uint      `json:"id" gorm:"primaryKey"`
    TitleEn          string    `json:"title_en"`
    TitleUa          string    `json:"title_ua"`
    DescriptionEn    string    `json:"description_en"`
    DescriptionUa    string    `json:"description_ua"`
    TimeTo           string    `json:"time_to"`
    TimeFrom         string    `json:"time_from"`
    AmountOfPersons  string    `json:"amount_of_persons"`
    ImageUrl         string    `json:"image_url"`
    CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
}
