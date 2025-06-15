package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Id       int    `json:"id" gorm:"primary_key"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-"`
}