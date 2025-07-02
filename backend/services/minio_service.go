package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/kholodihor/cows-shelter-backend/config"
	"github.com/minio/minio-go/v7"
)

// Ensure MinioService implements StorageService
var _ StorageService = (*MinioService)(nil)

// MinioService handles file operations with MinIO
type MinioService struct {
	client     *minio.Client
	bucketName string
	endpoint   string
	useSSL     bool
}

// NewMinioService creates a new MinioService instance
func NewMinioService() (*MinioService, error) {
	if config.MinioClient == nil {
		return nil, fmt.Errorf("MinIO client is not initialized")
	}

	bucketName := os.Getenv("MINIO_BUCKET")
	if bucketName == "" {
		bucketName = os.Getenv("STORAGE_BUCKET")
		if bucketName == "" {
			bucketName = "cows-shelter"
		}
	}

	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		endpoint = "minio:9000"
	}

	useSSL := os.Getenv("MINIO_USE_SSL") != "false"
	if useSSL {
		if os.Getenv("MINIO_USE_SSL") == "" && os.Getenv("STORAGE_USE_SSL") == "false" {
			useSSL = false
		}
	}

	return &MinioService{
		client:     config.MinioClient,
		bucketName: bucketName,
		endpoint:   endpoint,
		useSSL:     useSSL,
	}, nil
}

// UploadFile uploads a file to MinIO and returns the URL
func (s *MinioService) UploadFile(ctx context.Context, fileHeader *multipart.FileHeader, folder string) (string, error) {
	// Open the uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Generate a unique filename
	ext := filepath.Ext(fileHeader.Filename)
	objectName := fmt.Sprintf("%s/%s%s", folder, uuid.New().String(), ext)

	// Upload the file to MinIO
	_, err = s.client.PutObject(ctx, s.bucketName, objectName, file, fileHeader.Size, minio.PutObjectOptions{
		ContentType: fileHeader.Header.Get("Content-Type"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to MinIO: %v", err)
	}

	// Return the direct URL for the uploaded file
	return s.GetObjectURL(objectName), nil
}

// UploadBase64 uploads a base64-encoded image to MinIO and returns the URL
func (s *MinioService) UploadBase64(ctx context.Context, base64Data, folder string) (string, error) {
	// Extract the content type and decode the base64 data
	parts := strings.Split(base64Data, ";base64,")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid base64 data format")
	}

	contentType := strings.TrimPrefix(parts[0], "data:")
	base64Image := parts[1]

	// Decode the base64 data
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64Image))

	// Generate a unique filename with appropriate extension
	ext := ".jpg" // Default to jpg
	if strings.Contains(contentType, "png") {
		ext = ".png"
	} else if strings.Contains(contentType, "gif") {
		ext = ".gif"
	} else if strings.Contains(contentType, "webp") {
		ext = ".webp"
	}

	objectName := fmt.Sprintf("%s/%s%s", folder, uuid.New().String(), ext)

	// Upload the file to MinIO
	_, err := s.client.PutObject(ctx, s.bucketName, objectName, reader, -1, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload base64 data to MinIO: %v", err)
	}

	// Return the direct URL for the uploaded file
	return s.GetObjectURL(objectName), nil
}

// DeleteFile deletes a file from MinIO
func (s *MinioService) DeleteFile(ctx context.Context, objectName string) error {
	// Extract the object name from the URL if needed
	if strings.Contains(objectName, "/") {
		parts := strings.Split(objectName, s.bucketName+"/")
		if len(parts) > 1 {
			objectName = parts[1]
		}
	}

	// Delete the object from MinIO
	err := s.client.RemoveObject(ctx, s.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file from MinIO: %v", err)
	}

	return nil
}

// GetObjectURL returns the public URL for an object
func (s *MinioService) GetObjectURL(objectKey string) string {
	// Handle cases where objectKey might be a full URL
	if strings.HasPrefix(objectKey, "http") {
		return objectKey
	}

	// Construct the URL based on the endpoint and bucket
	var baseURL string
	if s.useSSL {
		baseURL = fmt.Sprintf("https://%s", s.endpoint)
	} else {
		baseURL = fmt.Sprintf("http://%s", s.endpoint)
	}

	// Return the direct URL for the object
	return fmt.Sprintf("%s/%s/%s", strings.TrimSuffix(baseURL, "/"), s.bucketName, strings.TrimPrefix(objectKey, "/"))
}

// ExtractObjectName extracts the object name from a MinIO URL
func (s *MinioService) ExtractObjectName(url string) string {
	// If it's not a full URL, return as is
	if !strings.Contains(url, "://") {
		return url
	}

	// Extract the object name from the URL
	parts := strings.Split(url, s.bucketName+"/")
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}
