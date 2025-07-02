package config

import (
	"fmt"
	"os"

	"strings"

	"github.com/kholodihor/cows-shelter-backend/storage"
	"github.com/kholodihor/cows-shelter-backend/storage/minio"
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
	// Default to S3 if no storage type is specified
	storageType := os.Getenv("STORAGE_TYPE")
	if storageType == "" {
		// Check for MinIO environment variables for backward compatibility
		if os.Getenv("MINIO_ENDPOINT") != "" {
			storageType = string(storage.TypeMinio)
		} else {
			storageType = string(storage.TypeS3)
		}
	}

	// Get common configuration
	bucketName := os.Getenv("STORAGE_BUCKET")
	if bucketName == "" {
		bucketName = os.Getenv("MINIO_BUCKET")
		if bucketName == "" {
			bucketName = "cows-shelter"
		}
	}

	useSSL := os.Getenv("STORAGE_USE_SSL") != "false"
	if useSSL {
		if os.Getenv("STORAGE_USE_SSL") == "" && os.Getenv("MINIO_USE_SSL") == "false" {
			useSSL = false
		}
	}

	// Get storage-specific configuration
	switch storage.Type(storageType) {
	case storage.TypeMinio:
		endpoint := os.Getenv("MINIO_ENDPOINT")
		if endpoint == "" {
			endpoint = "minio:9000"
		}

		// Ensure endpoint doesn't have a protocol
		endpoint = strings.TrimPrefix(endpoint, "http://")
		endpoint = strings.TrimPrefix(endpoint, "https://")

		return &Config{
			Type:       storage.TypeMinio,
			Endpoint:   endpoint,
			BucketName: bucketName,
			UseSSL:     useSSL,
		}

	case storage.TypeS3:
		endpoint := os.Getenv("S3_ENDPOINT")
		if endpoint == "" && os.Getenv("MINIO_ENDPOINT") != "" {
			// Fallback to MinIO endpoint if S3_ENDPOINT is not set but MINIO_ENDPOINT is
			endpoint = os.Getenv("MINIO_ENDPOINT")
		}

		region := os.Getenv("AWS_REGION")
		if region == "" {
			region = os.Getenv("MINIO_REGION")
			if region == "" {
				region = "us-east-1"
			}
		}

		return &Config{
			Type:       storage.TypeS3,
			Endpoint:   endpoint,
			Region:     region,
			BucketName: bucketName,
			UseSSL:     useSSL,
		}

	default:
		// Default to S3 configuration
		return &Config{
			Type:       storage.TypeS3,
			Region:     os.Getenv("AWS_REGION"),
			BucketName: bucketName,
			UseSSL:     true,
		}
	}
}

// NewStorageService creates a new storage service based on the configuration
func NewStorageService() (storage.Service, error) {
	cfg := GetConfig()

	switch cfg.Type {
	case storage.TypeMinio:
		// Initialize MinIO service
		return minio.New(
			cfg.Endpoint,
			cfg.BucketName,
			os.Getenv("MINIO_ACCESS_KEY"),
			os.Getenv("MINIO_SECRET_KEY"),
			cfg.UseSSL,
		)

	case storage.TypeS3:
		// Initialize S3 service
		return s3.New(
			cfg.Endpoint,
			cfg.BucketName,
			cfg.Region,
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			cfg.UseSSL,
		)

	default:
		return nil, fmt.Errorf("unsupported storage type: %s", cfg.Type)
	}
}
