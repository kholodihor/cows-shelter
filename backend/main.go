package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kholodihor/cows-shelter-backend/config"
	"github.com/kholodihor/cows-shelter-backend/handler"
	"github.com/kholodihor/cows-shelter-backend/middleware"
	"github.com/kholodihor/cows-shelter-backend/models"
	"golang.org/x/crypto/bcrypt"
)

// runMigrations runs database migrations directly without requiring Go toolchain
func runMigrations() error {
	// Initialize database connection
	config.Connect()
	
	// Run migrations
	if err := config.DB.AutoMigrate(
		&models.User{},
		&models.News{},
		&models.Excursion{},
		&models.Partner{},
		&models.Gallery{},
		&models.Review{},
		&models.Pdf{},
	); err != nil {
		return fmt.Errorf("failed to run migrations: %v", err)
	}
	
	log.Println("Migrations completed successfully")
	return nil
}

// createAdminUser creates or updates an admin user with the given email and password
func createAdminUser(email, password string) error {
	// Initialize database connection
	config.Connect()
	
	// Check if admin user already exists
	var user models.User
	result := config.DB.Where("email = ?", email).First(&user)
	
	hasher := bcrypt.DefaultCost
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), hasher)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	userData := models.User{
		Email:    email,
		Password: string(hashedPassword),
		Role:     "admin",
	}

	if result.Error != nil {
		// User doesn't exist, create new admin
		if err := config.DB.Create(&userData).Error; err != nil {
			return fmt.Errorf("failed to create admin user: %v", err)
		}
		log.Printf("Created new admin user: %s", email)
	} else {
		// Update existing user to admin
		user.Role = "admin"
		user.Password = string(hashedPassword)
		if err := config.DB.Save(&user).Error; err != nil {
			return fmt.Errorf("failed to update user to admin: %v", err)
		}
		log.Printf("Updated user to admin: %s", email)
	}

	return nil
}

func main() {
	// Handle command-line flags
	createAdminCmd := flag.NewFlagSet("create-admin", flag.ExitOnError)
	adminEmail := createAdminCmd.String("email", "", "Admin user email")
	adminPassword := createAdminCmd.String("password", "", "Admin user password")

	// Check if migrate command is provided
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		log.Println("Running database migrations...")
		if err := runMigrations(); err != nil {
			log.Fatalf("Failed to run migrations: %v\n", err)
		}
		log.Println("Migrations completed successfully")
		return
	}

	// Handle create-admin command
	if len(os.Args) > 1 && os.Args[1] == "create-admin" {
		if err := createAdminCmd.Parse(os.Args[2:]); err != nil {
			log.Fatalf("Failed to parse flags: %v\n", err)
		}

		if *adminEmail == "" || *adminPassword == "" {
			log.Fatal("Both --email and --password are required")
		}

		if err := createAdminUser(*adminEmail, *adminPassword); err != nil {
			log.Fatalf("Failed to create admin user: %v\n", err)
		}

		log.Printf("Successfully created/updated admin user: %s\n", *adminEmail)
		return
	}

	log.Println("Starting server...")

	router := gin.Default()

	config.Connect()

	// Initialize MinIO client
	if err := config.InitMinio(); err != nil {
		log.Printf("Warning: Failed to initialize MinIO client: %v\n", err)
	} else {
		log.Println("MinIO client initialized successfully")
	}

	// Initialize storage service
	storageService, err := config.NewStorageService()
	if err != nil {
		log.Printf("Warning: Failed to initialize storage service: %v\n", err)
	} else {
		log.Printf("Storage service initialized successfully (type: %s)", storageService)
		// Add storage middleware only if storage service was initialized successfully
		router.Use(middleware.GinStorageMiddleware(storageService))
	}

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		// Add security headers
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")

		c.Next()
	})

	// Initialize routes with the router
	handler.RouterHandler(&handler.Config{
		R: router,
	})

	router.Static("/uploads", "./uploads")

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to initialize server: %v\n", err)
		}
	}()

	log.Printf("Listening on port: %v\n", srv.Addr)

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Shutting down the server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}
}
