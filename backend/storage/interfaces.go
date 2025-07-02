package storage

import (
	"context"
	"mime/multipart"
	"time"
)

// ObjectInfo contains information about a stored object
type ObjectInfo struct {
	Key          string
	LastModified time.Time
	Size         int64
	ContentType  string
}

// Service defines the interface for storage operations
type Service interface {
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

	// ListObjects lists objects in the storage with the given prefix
	ListObjects(ctx context.Context, prefix string, maxKeys int) ([]ObjectInfo, error)
}

// Type represents the type of storage service
type Type string

const (
	// TypeMinio represents MinIO storage
	TypeMinio Type = "minio"
	// TypeS3 represents AWS S3 storage
	TypeS3 Type = "s3"
)
