package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadImage(c *gin.Context) {
	// Parse the file from the request
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload image: " + err.Error()})
		return
	}

	// Create an upload directory if it doesn't exist
	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err = os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create upload directory"})
			return
		}
	}

	// Save the file with a timestamp to avoid name conflicts
	filename := time.Now().Format("20060102150405") + "_" + filepath.Base(file.Filename)
	filepath := filepath.Join(uploadDir, filename)
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save file"})
		return
	}

	// Construct the image URL
	imageUrl := "/uploads/" + filename

	// Return the image URL in the response
	c.JSON(http.StatusCreated, gin.H{"image_url": imageUrl})
}