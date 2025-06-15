package utils

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
)

// AllowedImageTypes is a map of allowed image MIME types
var AllowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

// ValidateImageFile validates if the uploaded file is an image and within size limits
func ValidateImageFile(fileHeader *multipart.FileHeader, maxSize int64) error {
	// Check file size
	if fileHeader.Size > maxSize {
		return fmt.Errorf("file size exceeds the maximum limit of %d bytes", maxSize)
	}

	// Check file type
	if !AllowedImageTypes[fileHeader.Header.Get("Content-Type")] {
		return fmt.Errorf("invalid file type. Only JPEG, PNG, GIF, and WebP are allowed")
	}

	// Check file extension
	ext := filepath.Ext(fileHeader.Filename)
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		// Valid extension
	default:
		return fmt.Errorf("invalid file extension. Only .jpg, .jpeg, .png, .gif, and .webp are allowed")
	}

	return nil
}

// GetFileExtension returns the file extension from a filename
func GetFileExtension(filename string) string {
	return filepath.Ext(filename)
}
