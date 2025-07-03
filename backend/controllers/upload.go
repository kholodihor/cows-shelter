package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kholodihor/cows-shelter-backend/services"
)

var s3Service *services.S3Service

func init() {
	var err error
	s3Service, err = services.NewS3Service()
	if err != nil {
		panic("Failed to initialize S3 service: " + err.Error())
	}
}

func UploadImage(c *gin.Context) {
	// Parse the file from the request
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload image: " + err.Error()})
		return
	}

	// Upload the file to S3
	imageURL, err := s3Service.UploadFile(context.Background(), file, "uploads")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file to S3: " + err.Error()})
		return
	}

	// Return the S3 URL in the response
	c.JSON(http.StatusCreated, gin.H{"image_url": imageURL})
}