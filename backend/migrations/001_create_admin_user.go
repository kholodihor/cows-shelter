//go:build ignore

package main

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"unique"`
	Password string
}

func CreateAdminUser(db *gorm.DB) error {
	// Check if admin already exists
	var count int64
	if err := db.Model(&User{}).Where("email = ?", "admin@cowshelter.com").Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		log.Println("Admin user already exists")
		return nil
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create admin user
	admin := User{
		Email:    "admin@cowshelter.com",
		Password: string(hashedPassword),
	}

	if err := db.Create(&admin).Error; err != nil {
		return err
	}

	log.Println("Admin user created successfully")
	return nil
}
