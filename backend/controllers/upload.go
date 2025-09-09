package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kholodihor/cows-shelter-backend/storage"
)

func UploadImage(c *gin.Context) {
	// Get storage service from middleware
	storageService, exists := c.Get("storage")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage service not available"})
		return
	}

	service, ok := storageService.(storage.Service)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid storage service"})
		return
	}

	// Parse the file from the request
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload image: " + err.Error()})
		return
	}

	// Upload the file using storage service
	imageURL, err := service.UploadFile(context.Background(), file, "uploads")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file: " + err.Error()})
		return
	}

	// Return the file URL in the response
	c.JSON(http.StatusCreated, gin.H{"image_url": imageURL})
}