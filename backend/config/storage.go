package config

import (
	"os"

	"github.com/kholodihor/cows-shelter-backend/storage"
	"github.com/kholodihor/cows-shelter-backend/storage/cloudinary"
)

// Config holds the storage configuration
type Config struct {
	Type      storage.Type
	CloudName string
}

// GetConfig reads the storage configuration from environment variables
func GetConfig() *Config {
	return &Config{
		Type:      storage.TypeCloudinary,
		CloudName: os.Getenv("CLOUDINARY_CLOUD_NAME"),
	}
}

// GetStorageService returns the Cloudinary storage service
func GetStorageService(config *Config) (storage.Service, error) {
	return cloudinary.NewService(
		config.CloudName,
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
}
