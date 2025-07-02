package models

import "gorm.io/gorm"

type Pdf struct {
    gorm.Model
    Title       string `json:"title"`
    DocumentUrl string `json:"document_url"`
}