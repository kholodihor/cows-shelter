package config

import (
	"context"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	MinIOClient *minio.Client
)

// MinIOConfig holds MinIO configuration
type MinIOConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
}

// GetMinIOConfig loads MinIO configuration from environment variables
func GetMinIOConfig() *MinIOConfig {
	return &MinIOConfig{
		Endpoint:        GetEnv("MINIO_ENDPOINT", "localhost:9000"),
		AccessKeyID:     GetEnv("MINIO_ACCESS_KEY", "minioadmin"),
		SecretAccessKey: GetEnv("MINIO_SECRET_KEY", "minioadmin"),
		UseSSL:          GetEnv("MINIO_USE_SSL", "false") == "true",
		BucketName:      GetEnv("MINIO_BUCKET_NAME", "cows-shelter-uploads"),
	}
}

// InitMinIO initializes MinIO client and creates bucket if it doesn't exist
func InitMinIO() error {
	config := GetMinIOConfig()

	// Initialize MinIO client
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize MinIO client: %w", err)
	}

	MinIOClient = client

	// Check if bucket exists, create if it doesn't
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, config.BucketName)
	if err != nil {
		return fmt.Errorf("failed to check if bucket exists: %w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, config.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		log.Printf("Successfully created bucket: %s", config.BucketName)
	} else {
		log.Printf("Bucket %s already exists", config.BucketName)
	}

	// Set bucket policy to public read for uploaded files
	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {"AWS": "*"},
				"Action": "s3:GetObject",
				"Resource": "arn:aws:s3:::%s/*"
			}
		]
	}`, config.BucketName)

	err = client.SetBucketPolicy(ctx, config.BucketName, policy)
	if err != nil {
		log.Printf("Warning: Failed to set bucket policy (this is normal for some MinIO setups): %v", err)
	}

	log.Printf("MinIO client initialized successfully. Endpoint: %s, Bucket: %s", config.Endpoint, config.BucketName)
	return nil
}
