package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryService struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinaryService() (*CloudinaryService, error) {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Cloudinary: %v", err)
	}

	return &CloudinaryService{cld: cld}, nil
}

// UploadFile uploads a file to Cloudinary and returns the secure URL
func (s *CloudinaryService) UploadFile(ctx context.Context, fileHeader *multipart.FileHeader, folder string) (string, error) {
	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Create a temporary file
	tempFile, err := os.CreateTemp("", filepath.Ext(fileHeader.Filename))
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Copy the file content to the temp file
	if _, err = io.Copy(tempFile, file); err != nil {
		return "", fmt.Errorf("failed to copy file content: %v", err)
	}

	// Upload the file to Cloudinary
	overwrite := true
	result, err := s.cld.Upload.Upload(
		ctx,
		tempFile.Name(),
		uploader.UploadParams{
			Folder:   folder,
			PublicID:  fileHeader.Filename,
			Overwrite: &overwrite,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to Cloudinary: %v", err)
	}

	return result.SecureURL, nil
}

// DeleteFile deletes a file from Cloudinary
func (s *CloudinaryService) DeleteFile(ctx context.Context, publicID string) error {
	_, err := s.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	return err
}

// UploadBase64 uploads a base64-encoded image to Cloudinary and returns the secure URL
func (s *CloudinaryService) UploadBase64(ctx context.Context, base64Data, folder string) (string, error) {
	// The base64 data should be in the format: "data:image/png;base64,iVBORw0KGgo..."
	// We need to extract the actual base64 string
	parts := strings.SplitN(base64Data, ",", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid base64 data format")
	}

	// Upload the base64 data to Cloudinary
	uploadResult, err := s.cld.Upload.Upload(
		ctx,
		base64Data,
		uploader.UploadParams{
			Folder: folder,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload base64 image: %v", err)
	}

	return uploadResult.SecureURL, nil
}

// ExtractPublicID extracts the public ID from a Cloudinary URL
func (s *CloudinaryService) ExtractPublicID(url string) string {
	if url == "" {
		return ""
	}

	// Example URL formats:
	// https://res.cloudinary.com/<cloud_name>/image/upload/v<version>/<public_id>.<format>
	// https://res.cloudinary.com/<cloud_name>/image/upload/<public_id>.<format>
	// We need to extract the <public_id> part

	// Split by "upload/" and take the last part
	parts := strings.Split(url, "upload/")
	if len(parts) < 2 {
		return ""
	}

	// Get the part after "upload/"
	path := parts[1]

	// Remove any version prefix (v1234567890/)
	if strings.HasPrefix(path, "v") {
		if slashIndex := strings.Index(path, "/"); slashIndex > 0 {
			path = path[slashIndex+1:]
		}
	}

	// Remove file extension
	lastDotIndex := strings.LastIndex(path, ".")
	if lastDotIndex > 0 {
		path = path[:lastDotIndex]
	}

	// Remove any transformation parameters
	if slashIndex := strings.Index(path, "/"); slashIndex > 0 {
		path = path[:slashIndex]
	}

	// Remove the folder prefix if present
	folder := "cows-shelter/gallery/"
	if strings.HasPrefix(path, folder) {
		path = strings.TrimPrefix(path, folder)
	}

	return path
}
