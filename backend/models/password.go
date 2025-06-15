package models

import "gorm.io/gorm"


type Password struct {
    gorm.Model
    ID               uint      `json:"id" gorm:"primaryKey"`
    Email            string    `json:"email"`
    Token            string    `json:"token" gorm:"unique"`
}