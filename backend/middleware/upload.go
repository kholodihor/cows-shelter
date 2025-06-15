package middleware

import (
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kholodihor/cows-shelter-backend/utils"
)

const (
	// 5MB max file size
	maxUploadSize = 5 << 20
)

// UploadFile handles file uploads and validates them
func UploadFile(fieldName, folder string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the multipart form
		if err := c.Request.ParseMultipartForm(maxUploadSize); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form"})
			return
		}

		// Get the file from the request
		file, fileHeader, err := c.Request.FormFile(fieldName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
			return
		}
		defer file.Close()

		// Validate the file
		if err := utils.ValidateImageFile(fileHeader, maxUploadSize); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Store file info in context for the next handler
		c.Set("file", file)
		c.Set("fileHeader", fileHeader)
		c.Set("folder", folder)

		c.Next()
	}
}

// GetUploadedFile retrieves the uploaded file from the context
func GetUploadedFile(c *gin.Context) (multipart.File, *multipart.FileHeader, string, bool) {
	file, exists := c.Get("file")
	if !exists {
		return nil, nil, "", false
	}

	fileHeader, exists := c.Get("fileHeader")
	if !exists {
		return nil, nil, "", false
	}

	folder, exists := c.Get("folder")
	if !exists {
		folder = ""
	}

	return file.(multipart.File), fileHeader.(*multipart.FileHeader), folder.(string), true
}
