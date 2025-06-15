package services_test

import (
	"context"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"

	"github.com/kholodihor/cows-shelter-backend/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractPublicID(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "URL with version",
			url:      "https://res.cloudinary.com/demo/image/upload/v1234567890/folder/image.jpg",
			expected: "image",
		},
		{
			name:     "URL without version",
			url:      "https://res.cloudinary.com/demo/image/upload/folder/image.jpg",
			expected: "image",
		},
		{
			name:     "URL with transformation",
			url:      "https://res.cloudinary.com/demo/image/upload/w_100,h_100/folder/image.jpg",
			expected: "image",
		},
		{
			name:     "Empty URL",
			url:      "",
			expected: "",
		},
		{
			name:     "Invalid URL",
			url:      "not-a-valid-url",
			expected: "",
		},
	}

	service, err := services.NewCloudinaryService()
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.ExtractPublicID(tt.url)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUploadAndDeleteFile(t *testing.T) {
	// Skip this test in short mode to avoid making actual API calls
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a test file
	tempFile, err := os.CreateTemp("", "test-*.png")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// Write some test data to the file
	_, err = tempFile.Write([]byte("test image content"))
	require.NoError(t, err)
	tempFile.Close()

	// Create a file header for the test file
	file, err := os.Open(tempFile.Name())
	require.NoError(t, err)
	defer file.Close()

	fileInfo, err := file.Stat()
	require.NoError(t, err)

	fileHeader := &multipart.FileHeader{
		Filename: filepath.Base(tempFile.Name()),
		Size:     fileInfo.Size(),
	}

	// Initialize the service
	service, err := services.NewCloudinaryService()
	require.NoError(t, err)

	// Test UploadFile
	ctx := context.Background()
	folder := "cows-shelter/test"

	// Upload the file
	imageURL, err := service.UploadFile(ctx, fileHeader, folder)
	require.NoError(t, err)
	require.NotEmpty(t, imageURL)

	// Extract the public ID
	publicID := service.ExtractPublicID(imageURL)
	require.NotEmpty(t, publicID)

	// Test DeleteFile
	err = service.DeleteFile(ctx, publicID)
	assert.NoError(t, err)
}
