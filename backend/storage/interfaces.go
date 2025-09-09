package storage

import (
	"context"
	"mime/multipart"
	"time"
)

// Type represents the storage type
type Type string

const (
	// TypeCloudinary represents Cloudinary storage
	TypeCloudinary Type = "cloudinary"
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
	// UploadFile uploads a file to the storage and returns the URL
	UploadFile(ctx context.Context, file *multipart.FileHeader, folder string) (string, error)
	
	// UploadBase64 uploads a base64 encoded file to the storage and returns the URL
	UploadBase64(ctx context.Context, base64Data, filename, folder string) (string, error)
	
	// DeleteFile deletes a file from the storage
	DeleteFile(ctx context.Context, objectKey string) error
	
	// GetObjectURL returns the URL for accessing an object
	GetObjectURL(objectKey string) string
	
	// ExtractObjectName extracts the object name from a URL
	ExtractObjectName(url string) string
	
	// ListObjects lists objects in the storage with the given prefix
	ListObjects(ctx context.Context, prefix string, maxKeys int) ([]Object, error)
}

// Object represents a storage object
type Object struct {
	Key          string
	Size         int64
	LastModified string
	URL          string
}
