package cloudinary

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
	"github.com/kholodihor/cows-shelter-backend/storage"
)

type Service struct {
	client *cloudinary.Cloudinary
}

func NewService(cloudName, apiKey, apiSecret string) (*Service, error) {
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Cloudinary: %w", err)
	}

	return &Service{
		client: cld,
	}, nil
}

func (s *Service) UploadFile(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
	// Open the file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := uuid.New().String() + ext

	// Upload to Cloudinary
	uploadResult, err := s.client.Upload.Upload(ctx, src, uploader.UploadParams{
		PublicID: fmt.Sprintf("%s/%s", folder, strings.TrimSuffix(filename, ext)),
		Folder:   folder,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to Cloudinary: %w", err)
	}

	return uploadResult.SecureURL, nil
}

func (s *Service) DeleteFile(ctx context.Context, fileURL string) error {
	// Extract public ID from URL
	publicID := extractPublicIDFromURL(fileURL)
	if publicID == "" {
		return fmt.Errorf("invalid Cloudinary URL")
	}

	_, err := s.client.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete from Cloudinary: %w", err)
	}

	return nil
}

// UploadBase64 uploads a base64 encoded file to Cloudinary
func (s *Service) UploadBase64(ctx context.Context, base64Data, filename, folder string) (string, error) {
	// Parse base64 data (remove data:image/type;base64, prefix if present)
	if strings.Contains(base64Data, ",") {
		parts := strings.Split(base64Data, ",")
		if len(parts) > 1 {
			base64Data = parts[1]
		}
	}

	// Decode base64 data
	imageData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 data: %w", err)
	}

	// Create a bytes reader from the decoded data
	reader := bytes.NewReader(imageData)

	// Generate unique filename if not provided
	if filename == "" {
		filename = uuid.New().String()
	}

	// Remove extension from filename for public ID
	ext := filepath.Ext(filename)
	publicID := strings.TrimSuffix(filename, ext)

	// Upload to Cloudinary using the bytes reader
	uploadResult, err := s.client.Upload.Upload(ctx, reader, uploader.UploadParams{
		PublicID: fmt.Sprintf("%s/%s", folder, publicID),
		Folder:   folder,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to Cloudinary: %w", err)
	}

	return uploadResult.SecureURL, nil
}

// GetObjectURL returns the URL for accessing an object (for Cloudinary, this is just the URL itself)
func (s *Service) GetObjectURL(objectKey string) string {
	return objectKey
}

// ExtractObjectName extracts the object name from a Cloudinary URL
func (s *Service) ExtractObjectName(url string) string {
	return extractPublicIDFromURL(url)
}

func (s *Service) ListObjects(ctx context.Context, prefix string, limit int) ([]storage.Object, error) {
	// Cloudinary doesn't have a direct list API like S3
	// This is a placeholder implementation
	return []storage.Object{}, nil
}

func extractPublicIDFromURL(url string) string {
	// Extract public ID from Cloudinary URL
	// Example: https://res.cloudinary.com/demo/image/upload/v1234567890/folder/filename.jpg
	parts := strings.Split(url, "/")
	if len(parts) < 7 {
		return ""
	}
	
	// Find the upload part
	uploadIndex := -1
	for i, part := range parts {
		if part == "upload" {
			uploadIndex = i
			break
		}
	}
	
	if uploadIndex == -1 || uploadIndex+2 >= len(parts) {
		return ""
	}
	
	// Skip version if present (starts with 'v')
	startIndex := uploadIndex + 1
	if strings.HasPrefix(parts[startIndex], "v") {
		startIndex++
	}
	
	// Join remaining parts and remove extension
	publicID := strings.Join(parts[startIndex:], "/")
	if lastDot := strings.LastIndex(publicID, "."); lastDot != -1 {
		publicID = publicID[:lastDot]
	}
	
	return publicID
}
