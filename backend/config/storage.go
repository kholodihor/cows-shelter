package config

import (
	"os"

	"github.com/kholodihor/cows-shelter-backend/storage"
	"github.com/kholodihor/cows-shelter-backend/storage/s3"
)

// Config holds the storage configuration
type Config struct {
	Type       storage.Type
	Endpoint   string
	Region     string
	BucketName string
	UseSSL     bool
}

// GetConfig loads the storage configuration from environment variables
func GetConfig() *Config {
	// Always use S3 storage

	// Get common configuration
	bucketName := os.Getenv("STORAGE_BUCKET")
	if bucketName == "" {
		bucketName = os.Getenv("S3_BUCKET_NAME")
		if bucketName == "" {
			bucketName = "cows-shelter-uploads"
		}
	}

	// Always use SSL for S3
	useSSL := true

	// Get S3 configuration
	endpoint := os.Getenv("S3_ENDPOINT")

	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}

	return &Config{
		Type:       storage.TypeS3,
		Endpoint:   endpoint,
		Region:     region,
		BucketName: bucketName,
		UseSSL:     useSSL,
	}
}

// NewStorageService creates a new storage service based on the configuration
func NewStorageService() (storage.Service, error) {
	cfg := GetConfig()

	// Initialize S3 service
	return s3.New(
		cfg.Endpoint,
		cfg.BucketName,
		cfg.Region,
		os.Getenv("AWS_ACCESS_KEY_ID"),
		os.Getenv("AWS_SECRET_ACCESS_KEY"),
		cfg.UseSSL,
	)
}
