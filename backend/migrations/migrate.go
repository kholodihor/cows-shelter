//go:build ignore

package main

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/kholodihor/cows-shelter-backend/config"
	"github.com/kholodihor/cows-shelter-backend/models"
)

func main() {
	// Initialize database connection
	config.Connect()
	
	// Auto-migrate all models
	err := config.DB.AutoMigrate(
		&models.User{},
		&models.Contact{},
		&models.Excursion{},
		&models.Gallery{},
		&models.News{},
		&models.Partner{},
		&models.Password{},
		&models.Pdf{},
		&models.Review{},
	)

	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Database migrations completed successfully")

	// Create admin user
	if err := createAdminUser(config.DB); err != nil {
		log.Fatalf("Failed to create admin user: %v", err)
	}

	log.Println("Migration process completed successfully")
}

// createAdminUser creates an admin user if it doesn't already exist
func createAdminUser(db *gorm.DB) error {
	// Check if admin already exists
	var count int64
	if err := db.Model(&models.User{}).Where("email = ?", "admin@cowshelter.com").Count(&count).Error; err != nil {
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
	admin := models.User{
		Email:    "admin@cowshelter.com",
		Password: string(hashedPassword),
	}

	if err := db.Create(&admin).Error; err != nil {
		return err
	}

	log.Println("Admin user created successfully")
	return nil
}
