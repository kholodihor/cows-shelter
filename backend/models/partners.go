package models

import "gorm.io/gorm"

type Partner struct {
    gorm.Model
    Name  string `json:"name"`
    Logo  string `json:"logo"`
    Link  string `json:"link"`
}
