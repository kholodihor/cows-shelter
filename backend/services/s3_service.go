package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// S3Service handles file operations with AWS S3
type S3Service struct {
	client       *s3.Client
	bucketName   string
	region       string
	endpoint     string
	usePathStyle bool
}

// NewS3Service creates a new S3Service instance with the provided configuration
func NewS3Service() (*S3Service, error) {
	// Get configuration from environment variables
	bucketName := os.Getenv("S3_BUCKET")
	if bucketName == "" {
		bucketName = "cows-shelter-uploads"
	}

	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1" // Default region
	}

	endpoint := os.Getenv("S3_ENDPOINT")

	useSSL := os.Getenv("S3_USE_SSL") != "false"
	if useSSL {
		if !strings.HasPrefix(endpoint, "http") {
			endpoint = "https://" + endpoint
		}
	} else if !strings.HasPrefix(endpoint, "http") {
		endpoint = "http://" + endpoint
	}

	// Configure AWS SDK with default credentials chain (will use ECS task role when running in ECS)
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %v", err)
	}

	// Log the credentials source for debugging
	creds, err := cfg.Credentials.Retrieve(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve AWS credentials: %v", err)
	}
	fmt.Printf("Using AWS credentials from: %s\n", creds.Source)

	// Create S3 client
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
			o.UsePathStyle = true // Use path style for custom endpoints
		}
	})

	return &S3Service{
		client:       client,
		bucketName:   bucketName,
		region:       region,
		endpoint:     endpoint,
		usePathStyle: true,
	}, nil
}

// UploadFile uploads a file to S3 and returns the URL
func (s *S3Service) UploadFile(ctx context.Context, fileHeader *multipart.FileHeader, folder string) (string, error) {
	// Open the uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Generate a unique filename
	ext := filepath.Ext(fileHeader.Filename)
	objectKey := fmt.Sprintf("%s/%s%s", folder, uuid.New().String(), ext)

	// Upload the file to S3
	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(objectKey),
		Body:        file,
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %v", err)
	}

	return s.GetObjectURL(objectKey), nil
}

// UploadBase64 uploads a base64-encoded image to S3 and returns the URL
func (s *S3Service) UploadBase64(ctx context.Context, base64Data, folder string) (string, error) {
	// Extract the content type and decode the base64 data
	parts := strings.SplitN(base64Data, ";base64,", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid base64 data format")
	}

	contentType := strings.TrimPrefix(parts[0], "data:")
	base64Image := parts[1]

	// Decode the base64 data
	imageData, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 data: %v", err)
	}

	// Generate a unique filename with appropriate extension
	ext := ".jpg" // Default to jpg
	if strings.Contains(contentType, "png") {
		ext = ".png"
	} else if strings.Contains(contentType, "gif") {
		ext = ".gif"
	} else if strings.Contains(contentType, "webp") {
		ext = ".webp"
	}

	objectKey := fmt.Sprintf("%s/%s%s", folder, uuid.New().String(), ext)

	// Upload the file to S3
	contentLength := int64(len(imageData))
	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucketName),
		Key:           aws.String(objectKey),
		Body:          bytes.NewReader(imageData),
		ContentType:   aws.String(contentType),
		ContentLength: &contentLength,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload base64 data to S3: %v", err)
	}

	return s.GetObjectURL(objectKey), nil
}

// DeleteFile deletes a file from S3
func (s *S3Service) DeleteFile(ctx context.Context, objectKey string) error {
	// Extract the object key from the URL if needed
	if strings.Contains(objectKey, s.bucketName) {
		parts := strings.SplitN(objectKey, s.bucketName+"/", 2)
		if len(parts) > 1 {
			objectKey = parts[1]
		}
	}

	// Delete the object from S3
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %v", err)
	}

	return nil
}

// GetObjectURL returns the public URL for an object
func (s *S3Service) GetObjectURL(objectKey string) string {
	// Handle cases where objectKey might be a full URL
	if strings.HasPrefix(objectKey, "http") {
		return objectKey
	}

	// Check if PUBLIC_STORAGE_URL environment variable is set (for CloudFront)
	publicStorageURL := os.Getenv("PUBLIC_STORAGE_URL")
	if publicStorageURL != "" {
		// Use CloudFront URL for public access
		return fmt.Sprintf("%s/%s", strings.TrimSuffix(publicStorageURL, "/"), strings.TrimPrefix(objectKey, "/"))
	}

	// For S3-compatible storage with custom endpoint
	if s.endpoint != "" {
		// Remove any protocol prefix if present
		endpoint := strings.TrimPrefix(s.endpoint, "http://")
		endpoint = strings.TrimPrefix(endpoint, "https://")
		return fmt.Sprintf("https://%s/%s/%s", strings.TrimSuffix(endpoint, "/"), s.bucketName, strings.TrimPrefix(objectKey, "/"))
	}

	// For AWS S3 with proper URL formatting
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucketName, s.region, strings.TrimPrefix(objectKey, "/"))
}

// ExtractObjectName extracts the object key from an S3 URL
func (s *S3Service) ExtractObjectName(url string) string {
	// Extract the object key from the URL
	parts := strings.Split(url, s.bucketName+"/")
	if len(parts) > 1 {
		return parts[1]
	}
	return url // Return as is if no bucket name found
}
