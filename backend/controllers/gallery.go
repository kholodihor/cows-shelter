package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kholodihor/cows-shelter-backend/config"
	"github.com/kholodihor/cows-shelter-backend/models"
	"github.com/kholodihor/cows-shelter-backend/services"
)

// CreateGalleryRequest represents the JSON request body for creating a gallery item
type CreateGalleryRequest struct {
	ImageData string `json:"image_data" binding:"required"` // base64-encoded image data
}

// UpdateGalleryRequest represents the JSON request body for updating a gallery item
type UpdateGalleryRequest struct {
	ImageData string `json:"image_data"` // base64-encoded image data (optional)
}

func GetAllGalleries(c *gin.Context) {
	galleries := []models.Gallery{}
	config.DB.Find(&galleries)
	c.JSON(http.StatusOK, &galleries)
}

func GetGalleries(c *gin.Context) {
	var galleries []models.Gallery
	var total int64

	// Default values for pagination
	limit := 10  // Default limit of 10 items per page
	page := 1    // Default to the first page

	// Parse limit and page from query parameters (if provided)
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}

	// Calculate the offset (skip items)
	offset := (page - 1) * limit

	// Count the total number of records
	config.DB.Model(&models.Gallery{}).Count(&total)

	// Fetch the paginated results
	if err := config.DB.Limit(limit).Offset(offset).Find(&galleries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching galleries"})
		return
	}

	// Return paginated response
	c.JSON(http.StatusOK, gin.H{
		"data":       galleries,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + int64(limit) - 1) / int64(limit), // Calculate total pages
	})
}

func GetGalleryByID(c *gin.Context) {
	var gallery models.Gallery
	id := c.Param("id")
	if err := config.DB.Where("id = ?", id).First(&gallery).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gallery not found"})
		return
	}
	c.JSON(http.StatusOK, &gallery)
}

// CreateGallery handles the creation of a new gallery item with a base64-encoded image
// @Summary Create a new gallery item
// @Description Create a new gallery item with a base64-encoded image
// @Tags gallery
// @Accept json
// @Produce json
// @Param input body CreateGalleryRequest true "Gallery item data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /gallery [post]
func CreateGallery(c *gin.Context) {
	var req CreateGalleryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	if req.ImageData == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image data is required"})
		return
	}

	// Initialize MinIO service
	minioService, err := services.NewMinioService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize MinIO service"})
		return
	}

	// Upload the base64 image to MinIO
	imageURL, err := minioService.UploadBase64(c.Request.Context(), req.ImageData, "gallery")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image: " + err.Error()})
		return
	}

	// Create the gallery entry in the database
	gallery := models.Gallery{
		ImageUrl:  imageURL,
		CreatedAt: time.Now(),
	}

	if err := config.DB.Create(&gallery).Error; err != nil {
		// If database save fails, try to delete the uploaded image from MinIO
		_ = minioService.DeleteFile(c.Request.Context(), imageURL)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create gallery entry: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":       gallery.ID,
		"imageUrl": gallery.ImageUrl,
	})
}

// UpdateGallery handles updating a gallery item with an optional new base64-encoded image
// @Summary Update a gallery item
// @Description Update a gallery item with an optional new base64-encoded image
// @Tags gallery
// @Accept json
// @Produce json
// @Param id path int true "Gallery ID"
// @Param input body UpdateGalleryRequest true "Updated gallery item data"
// @Success 200 {object} models.Gallery
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /gallery/{id} [put]
func UpdateGallery(c *gin.Context) {
	// Get the existing gallery item
	var gallery models.Gallery
	if err := config.DB.Where("id = ?", c.Param("id")).First(&gallery).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gallery not found"})
		return
	}

	var req UpdateGalleryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	// Initialize MinIO service
	minioService, err := services.NewMinioService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize MinIO service"})
		return
	}

	// Check if a new image is being uploaded
	if req.ImageData != "" {
		// Upload the new base64 image to MinIO
		newImageURL, err := minioService.UploadBase64(c.Request.Context(), req.ImageData, "gallery")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload new image: " + err.Error()})
			return
		}

		// Delete the old image from MinIO if it exists
		if gallery.ImageUrl != "" {
			_ = minioService.DeleteFile(c.Request.Context(), gallery.ImageUrl)
		}

		// Update the image URL
		gallery.ImageUrl = newImageURL

		// Save the updated gallery item
		if err := config.DB.Save(&gallery).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update gallery: " + err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, &gallery)
}

// DeleteGallery handles the deletion of a gallery item and its associated image
// @Summary Delete a gallery item
// @Description Delete a gallery item and its associated image
// @Tags gallery
// @Produce json
// @Param id path int true "Gallery ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /gallery/{id} [delete]
func DeleteGallery(c *gin.Context) {
	// Get the gallery item to delete
	var gallery models.Gallery
	if err := config.DB.Where("id = ?", c.Param("id")).First(&gallery).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gallery not found"})
		return
	}

	// Initialize MinIO service
	minioService, err := services.NewMinioService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize MinIO service"})
		return
	}

	// Delete the image from MinIO if it exists
	if gallery.ImageUrl != "" {
		// Extract the object name from the URL
		objectName := minioService.ExtractObjectName(gallery.ImageUrl)
		if objectName != "" {
			_ = minioService.DeleteFile(c.Request.Context(), objectName)
		}
	}

	// Delete the gallery item from the database
	if err := config.DB.Delete(&gallery).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete gallery"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Gallery item deleted successfully"})
}
