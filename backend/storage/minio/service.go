package minio

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/kholodihor/cows-shelter-backend/storage"
)

// Service implements the storage.Service interface for MinIO
type Service struct {
	client     *minio.Client
	bucketName string
	endpoint   string
	useSSL     bool
}

// New creates a new MinIO storage service
func New(endpoint, bucketName, accessKey, secretKey string, useSSL bool) (*Service, error) {
	// Initialize MinIO client
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	// Create bucket if it doesn't exist
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check if bucket exists: %w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}

		// Set bucket policy to public read
		policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": "*",
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::%s/*"]
				}
			]
		}`, bucketName)

		err = client.SetBucketPolicy(ctx, bucketName, policy)
		if err != nil {
			return nil, fmt.Errorf("failed to set bucket policy: %w", err)
		}
	}

	return &Service{
		client:     client,
		bucketName: bucketName,
		endpoint:   endpoint,
		useSSL:     useSSL,
	}, nil
}

// UploadFile uploads a file from a multipart file header
func (s *Service) UploadFile(ctx context.Context, fileHeader *multipart.FileHeader, folder string) (string, error) {
	// Open the uploaded file
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Generate a unique filename
	ext := filepath.Ext(fileHeader.Filename)
	objectName := generateObjectName(folder, ext)

	// Upload the file
	_, err = s.client.PutObject(
		ctx,
		s.bucketName,
		objectName,
		src,
		fileHeader.Size,
		minio.PutObjectOptions{
			ContentType: fileHeader.Header.Get("Content-Type"),
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to MinIO: %w", err)
	}

	return s.GetObjectURL(objectName), nil
}

// UploadBase64 uploads a base64-encoded image
func (s *Service) UploadBase64(ctx context.Context, base64Data, folder string) (string, error) {
	// Extract content type and data from base64 string
	contentType, data, err := parseBase64Data(base64Data)
	if err != nil {
		return "", fmt.Errorf("invalid base64 data: %w", err)
	}

	// Generate a unique filename
	ext := "." + strings.Split(contentType, "/")[1] // e.g., "image/png" -> ".png"
	objectName := generateObjectName(folder, ext)

	// Upload the file
	_, err = s.client.PutObject(
		ctx,
		s.bucketName,
		objectName,
		strings.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload base64 data to MinIO: %w", err)
	}

	return s.GetObjectURL(objectName), nil
}

// DeleteFile deletes a file from MinIO
func (s *Service) DeleteFile(ctx context.Context, objectName string) error {
	// If the object name is a full URL, extract just the object key
	objectKey := s.ExtractObjectName(objectName)

	// Delete the object
	err := s.client.RemoveObject(ctx, s.bucketName, objectKey, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object from MinIO: %w", err)
	}

	return nil
}

// GetObjectURL returns the public URL for an object
func (s *Service) GetObjectURL(objectKey string) string {
	// Handle cases where objectKey might be a full URL
	if strings.HasPrefix(objectKey, "http") {
		return objectKey
	}

	// Construct the URL based on the endpoint and bucket
	var scheme string
	if s.useSSL {
		scheme = "https://"
	} else {
		scheme = "http://"
	}

	// If using MinIO with a path-style URL
	return fmt.Sprintf("%s%s/%s/%s", scheme, s.endpoint, s.bucketName, strings.TrimPrefix(objectKey, "/"))
}

// ListObjects lists objects in the storage with the given prefix
func (s *Service) ListObjects(ctx context.Context, prefix string, maxKeys int) ([]storage.ObjectInfo, error) {
	// Ensure maxKeys is at least 1 and at most 1000
	if maxKeys < 1 {
		maxKeys = 1000
	}

	// List objects with the given prefix
	objectCh := s.client.ListObjects(ctx, s.bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: false,
		MaxKeys:   maxKeys,
	})

	var objects []storage.ObjectInfo
	for obj := range objectCh {
		if obj.Err != nil {
			return nil, fmt.Errorf("error listing objects: %w", obj.Err)
		}

		objects = append(objects, storage.ObjectInfo{
			Key:          obj.Key,
			LastModified: obj.LastModified,
			Size:         obj.Size,
			ContentType:  obj.ContentType,
		})

		// Break if we've reached the maximum number of keys
		if len(objects) >= maxKeys {
			break
		}
	}

	return objects, nil
}

// ExtractObjectName extracts the object name from a MinIO URL
func (s *Service) ExtractObjectName(url string) string {
	// If it's not a full URL, return as is
	if !strings.Contains(url, "://") {
		return url
	}

	// Extract the object name from the URL
	parts := strings.Split(url, s.bucketName+"/")
	if len(parts) > 1 {
		return parts[1]
	}

	// If we couldn't extract the object name, return the original URL
	return url
}

// generateObjectName generates a unique object name with the given folder and extension
func generateObjectName(folder, ext string) string {
	// Generate a unique filename
	filename := uuid.New().String() + ext

	// If folder is empty, just return the filename
	if folder == "" {
		return filename
	}

	// Otherwise, join the folder and filename
	return strings.TrimSuffix(folder, "/") + "/" + filename
}

// parseBase64Data parses a base64-encoded data URL and returns the content type and data
func parseBase64Data(dataURL string) (string, string, error) {
	// Format: data:image/png;base64,<data>
	parts := strings.SplitN(dataURL, ",", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid base64 data URL")
	}

	header := parts[0]
	if !strings.HasSuffix(header, "base64") {
		return "", "", fmt.Errorf("only base64 encoding is supported")
	}

	// Extract content type
	contentType := strings.TrimPrefix(header, "data:")
	contentType = strings.TrimSuffix(contentType, ";base64")

	// Return the base64 data as-is, let the caller decode it
	return contentType, parts[1], nil
}
