package s3

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
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"

	"github.com/kholodihor/cows-shelter-backend/storage"
)

// Service implements the storage.Service interface for AWS S3
type Service struct {
	client     *s3.Client
	bucketName string
	endpoint   string
	region     string
	useSSL     bool
}

// New creates a new S3 storage service
func New(endpoint, bucketName, region, accessKey, secretKey string, useSSL bool) (*Service, error) {
	// Create custom resolver for S3 endpoint
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		// Only use custom endpoint if provided, otherwise use AWS default
		if endpoint != "" {
			return aws.Endpoint{
				URL:               endpoint,
				SigningRegion:     region,
				HostnameImmutable: true,
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	// Create AWS config with credentials and custom endpoint
	cfg, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithRegion(region),
		awsconfig.WithEndpointResolverWithOptions(customResolver),
	)

	// If we have explicit credentials, use them (for MinIO or custom S3)
	if accessKey != "" && secretKey != "" {
		cfg.Credentials = credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true // Required for MinIO and some S3-compatible services
	})

		// Verify the bucket exists and is accessible
	ctx := context.Background()
	_, err = client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		return nil, fmt.Errorf("bucket %s does not exist or is not accessible: %w. Ensure the bucket exists and the IAM role has the necessary permissions", bucketName, err)
	}

	return &Service{
		client:     client,
		bucketName: bucketName,
		endpoint:   endpoint,
		region:     region,
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
	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(objectName),
		Body:        src,
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return s.GetObjectURL(objectName), nil
}

// UploadBase64 uploads a base64-encoded image
func (s *Service) UploadBase64(ctx context.Context, base64Data, folder string) (string, error) {
	// Extract content type and base64 data from the data URL
	contentType, base64Image, err := parseBase64Data(base64Data)
	if err != nil {
		return "", fmt.Errorf("invalid base64 data: %w", err)
	}

	// Decode the base64 data to binary
	imageData, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 data: %w", err)
	}

	// Generate a unique filename
	ext := "." + strings.Split(contentType, "/")[1] // e.g., "image/png" -> ".png"
	objectName := generateObjectName(folder, ext)

	// Upload the binary data to S3
	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(objectName),
		Body:        bytes.NewReader(imageData), // Use the decoded binary data
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload base64 data to S3: %w", err)
	}

	return s.GetObjectURL(objectName), nil
}

// DeleteFile deletes a file from S3
func (s *Service) DeleteFile(ctx context.Context, objectName string) error {
	// If the object name is a full URL, extract just the object key
	objectKey := s.ExtractObjectName(objectName)

	// Delete the object
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(objectKey),
	})

	if err != nil {
		return fmt.Errorf("failed to delete object from S3: %w", err)
	}

	return nil
}

// GetObjectURL returns the public URL for an object
func (s *Service) GetObjectURL(objectKey string) string {
	// Handle cases where objectKey might be a full URL
	if strings.HasPrefix(objectKey, "http") {
		return objectKey
	}

	// Check for public storage URL from environment variable
	if publicURL := os.Getenv("PUBLIC_STORAGE_URL"); publicURL != "" {
		// Ensure the public URL ends with a single trailing slash
		baseURL := strings.TrimSuffix(publicURL, "/")
		return fmt.Sprintf("%s/%s", baseURL, strings.TrimPrefix(objectKey, "/"))
	}

	// If using a custom endpoint (e.g., MinIO), construct the URL directly
	if s.endpoint != "" {
		scheme := "http://"
		if s.useSSL {
			scheme = "https://"
		}
		return fmt.Sprintf("%s%s/%s/%s", scheme, s.endpoint, s.bucketName, strings.TrimPrefix(objectKey, "/"))
	}

	// For AWS S3, use the standard URL format
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucketName, s.region, strings.TrimPrefix(objectKey, "/"))
}

// ListObjects lists objects in the storage with the given prefix
func (s *Service) ListObjects(ctx context.Context, prefix string, maxKeys int) ([]storage.ObjectInfo, error) {
	// Ensure maxKeys is at least 1 and at most 1000
	if maxKeys < 1 {
		maxKeys = 1
	} else if maxKeys > 1000 {
		maxKeys = 1000
	}

	// List objects with the given prefix
	maxKeysInt32 := int32(maxKeys)
	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(s.bucketName),
		Prefix:  aws.String(prefix),
		MaxKeys: &maxKeysInt32,
	}

	result, err := s.client.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %w", err)
	}

	var objects []storage.ObjectInfo
	for _, obj := range result.Contents {
		objects = append(objects, storage.ObjectInfo{
			Key:          aws.ToString(obj.Key),
			LastModified: aws.ToTime(obj.LastModified),
			Size:         aws.ToInt64(obj.Size),
			// ContentType is not directly available in ListObjectsV2 response
			// It would need to be fetched via HeadObject if required
			ContentType:  "",
		})

		// Break if we've reached the maximum number of keys
		if len(objects) >= maxKeys {
			break
		}
	}

	return objects, nil
}

// ExtractObjectName extracts the object name from an S3 URL
func (s *Service) ExtractObjectName(url string) string {
	// Remove the protocol and domain parts
	parts := strings.Split(url, "//")
	if len(parts) < 2 {
		return url
	}

	// Remove the bucket name and any path
	pathParts := strings.SplitN(parts[1], "/", 2)
	if len(pathParts) < 2 {
		return ""
	}

	// Handle custom endpoint URL format: http://endpoint/bucket/key
	endpointParts := strings.SplitN(url, s.bucketName+"/", 2)
	if len(endpointParts) > 1 {
		return endpointParts[1]
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
	// Handle data URL format: data:image/png;base64,iVBORw0KGgo...
	parts := strings.SplitN(dataURL, ",", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid base64 data URL format")
	}

	// Extract content type
	header := parts[0]
	if !strings.HasPrefix(header, "data:") || !strings.Contains(header, ";base64") {
		return "", "", fmt.Errorf("invalid base64 data URL format")
	}

	contentType := strings.TrimPrefix(header, "data:")
	contentType = strings.TrimSuffix(contentType, ";base64")

	return contentType, parts[1], nil
}
