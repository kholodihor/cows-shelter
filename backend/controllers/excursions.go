package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kholodihor/cows-shelter-backend/config"
	"github.com/kholodihor/cows-shelter-backend/middleware"
	"github.com/kholodihor/cows-shelter-backend/models"
)

// CreateExcursionRequest represents the JSON request body for creating an excursion
type CreateExcursionRequest struct {
	TitleEn         string `json:"title_en" binding:"required"`
	TitleUa         string `json:"title_ua"`
	DescriptionEn   string `json:"description_en"`
	DescriptionUa   string `json:"description_ua"`
	TimeTo          string `json:"time_to" binding:"required"`
	TimeFrom        string `json:"time_from" binding:"required"`
	AmountOfPersons string `json:"amount_of_persons" binding:"required"`
	ImageData       string `json:"image_data" binding:"required"` // base64-encoded image data
}

// UpdateExcursionRequest represents the JSON request body for updating an excursion
type UpdateExcursionRequest struct {
	TitleEn         string `json:"title_en"`
	TitleUa         string `json:"title_ua"`
	DescriptionEn   string `json:"description_en"`
	DescriptionUa   string `json:"description_ua"`
	TimeTo          string `json:"time_to"`
	TimeFrom        string `json:"time_from"`
	AmountOfPersons string `json:"amount_of_persons"`
	ImageData       string `json:"image_data"` // base64-encoded image data (optional)
}

func GetAllExcursions(c *gin.Context) {
	excursions := []models.Excursion{}
	config.DB.Find(&excursions)
	c.JSON(http.StatusOK, &excursions)
}

func GetExcursions(c *gin.Context) {
	var excursions []models.Excursion
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
	if err := config.DB.Model(&models.Excursion{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error counting excursions"})
		return
	}

	// Fetch the paginated results
	if err := config.DB.Limit(limit).Offset(offset).Find(&excursions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching excursions"})
		return
	}

	// Return paginated response
	c.JSON(http.StatusOK, gin.H{
		"excursions":   excursions,
		"totalLength":  total, // Total number of excursions in the database
		"page":         page,
		"limit":        limit,
		"totalPages":   (total + int64(limit) - 1) / int64(limit), // Calculate total pages
	})
}

func GetExcursionByID(c *gin.Context) {
	var excursion models.Excursion
	id := c.Param("id")
	if err := config.DB.Where("id = ?", id).First(&excursion).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Excursion not found"})
		return
	}
	c.JSON(http.StatusOK, &excursion)
}

// CreateExcursion handles the creation of a new excursion with a base64-encoded image
// @Summary Create a new excursion
// @Description Create a new excursion with a base64-encoded image
// @Tags excursions
// @Accept json
// @Produce json
// @Param input body CreateExcursionRequest true "Excursion data"
// @Success 201 {object} models.Excursion
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /excursions [post]
func CreateExcursion(c *gin.Context) {
	var req CreateExcursionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	// Create excursion model
	excursion := models.Excursion{
		TitleEn:         req.TitleEn,
		TitleUa:         req.TitleUa,
		DescriptionEn:   req.DescriptionEn,
		DescriptionUa:   req.DescriptionUa,
		TimeTo:          req.TimeTo,
		TimeFrom:        req.TimeFrom,
		AmountOfPersons: req.AmountOfPersons,
	}

	// Handle image upload
	if req.ImageData == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image data is required"})
		return
	}

	// Get storage service from context
	store := middleware.GetStorage(c.Request.Context())
	if store == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage service not available"})
		return
	}

	// Upload the base64 image using the storage service
	imageURL, err := store.UploadBase64(c.Request.Context(), req.ImageData, "excursions")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image: " + err.Error()})
		return
	}
	excursion.ImageUrl = imageURL

	// Save the excursion to the database
	if err := config.DB.Create(&excursion).Error; err != nil {
		// If database save fails, try to delete the uploaded image
		_ = store.DeleteFile(c.Request.Context(), imageURL)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create excursion: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, excursion)
}

// UpdateExcursion handles updating an existing excursion with an optional new base64-encoded image
// @Summary Update an existing excursion
// @Description Update an existing excursion with an optional new base64-encoded image
// @Tags excursions
// @Accept json
// @Produce json
// @Param id path int true "Excursion ID"
// @Param input body UpdateExcursionRequest true "Updated excursion data"
// @Success 200 {object} models.Excursion
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /excursions/{id} [put]
func UpdateExcursion(c *gin.Context) {
	// Get the existing excursion
	var excursion models.Excursion
	if err := config.DB.Where("id = ?", c.Param("id")).First(&excursion).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Excursion not found"})
		return
	}

	var req UpdateExcursionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	// Update fields if provided
	if req.TitleEn != "" {
		excursion.TitleEn = req.TitleEn
	}
	if req.TitleUa != "" {
		excursion.TitleUa = req.TitleUa
	}
	if req.DescriptionEn != "" {
		excursion.DescriptionEn = req.DescriptionEn
	}
	if req.DescriptionUa != "" {
		excursion.DescriptionUa = req.DescriptionUa
	}
	if req.TimeTo != "" {
		excursion.TimeTo = req.TimeTo
	}
	if req.TimeFrom != "" {
		excursion.TimeFrom = req.TimeFrom
	}
	if req.AmountOfPersons != "" {
		excursion.AmountOfPersons = req.AmountOfPersons
	}

	// Handle image upload if new image data is provided
	if req.ImageData != "" {
		// Get storage service from context
		store := middleware.GetStorage(c.Request.Context())
		if store == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage service not available"})
			return
		}

		// Upload the new base64 image using the storage service
		newImageURL, err := store.UploadBase64(c.Request.Context(), req.ImageData, "excursions")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload new image: " + err.Error()})
			return
		}

		// Delete the old image if it exists
		if excursion.ImageUrl != "" {
			oldObjectName := store.ExtractObjectName(excursion.ImageUrl)
			if oldObjectName != "" {
				_ = store.DeleteFile(c.Request.Context(), oldObjectName)
			}
		}

		// Update the image URL
		excursion.ImageUrl = newImageURL
	}

	// Save the updated excursion to the database
	if err := config.DB.Save(&excursion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update excursion: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, excursion)
}

// DeleteExcursion handles the deletion of an excursion and its associated image
// @Summary Delete an excursion
// @Description Delete an excursion and its associated image
// @Tags excursions
// @Produce json
// @Param id path int true "Excursion ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /excursions/{id} [delete]
func DeleteExcursion(c *gin.Context) {
	// Get the excursion to delete
	var excursion models.Excursion
	if err := config.DB.Where("id = ?", c.Param("id")).First(&excursion).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Excursion not found"})
		return
	}

	// Get storage service from context
	store := middleware.GetStorage(c.Request.Context())
	if store == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage service not available"})
		return
	}

	// Delete the image if it exists
	if excursion.ImageUrl != "" {
		// Extract the object name from the URL
		objectName := store.ExtractObjectName(excursion.ImageUrl)
		if objectName != "" {
			_ = store.DeleteFile(c.Request.Context(), objectName)
		}
	}

	// Delete the excursion from the database
	if err := config.DB.Delete(&excursion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete excursion"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Excursion deleted successfully"})
}