package services

import (
	"context"
	"mime/multipart"
)

// StorageService defines the interface for storage operations
type StorageService interface {
	// UploadFile uploads a file from a multipart file header and returns the URL
	UploadFile(ctx context.Context, fileHeader *multipart.FileHeader, folder string) (string, error)

	// UploadBase64 uploads a base64-encoded image and returns the URL
	UploadBase64(ctx context.Context, base64Data, folder string) (string, error)

	// DeleteFile deletes a file from storage
	DeleteFile(ctx context.Context, objectName string) error

	// GetObjectURL returns the public URL for an object
	GetObjectURL(objectKey string) string

	// ExtractObjectName extracts the object name/key from a URL
	ExtractObjectName(url string) string
}

// StorageType represents the type of storage service
type StorageType string

const (
	// StorageTypeMinio represents MinIO storage
	StorageTypeMinio StorageType = "minio"
	// StorageTypeS3 represents AWS S3 storage
	StorageTypeS3 StorageType = "s3"
)

// NewStorageService creates a new storage service based on the configuration
func NewStorageService(storageType StorageType) (StorageService, error) {
	switch storageType {
	case StorageTypeMinio:
		return NewMinioService()
	case StorageTypeS3:
		return NewS3Service()
	default:
		// Default to S3
		return NewS3Service()
	}
}
