package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kholodihor/cows-shelter-backend/config"
	"github.com/kholodihor/cows-shelter-backend/handler"
)

// runMigrations runs database migrations
func runMigrations() error {
	cmd := exec.Command("go", "run", "migrations/migrate.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	// Check if migrate command is provided
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		log.Println("Running database migrations...")
		if err := runMigrations(); err != nil {
			log.Fatalf("Failed to run migrations: %v\n", err)
		}
		log.Println("Migrations completed successfully")
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
