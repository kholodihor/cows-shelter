package models

import "gorm.io/gorm"


type Password struct {
    gorm.Model
    Email string `json:"email"`
    Token string `json:"token" gorm:"unique"`
}