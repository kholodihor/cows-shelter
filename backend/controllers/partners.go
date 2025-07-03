package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kholodihor/cows-shelter-backend/config"
	"github.com/kholodihor/cows-shelter-backend/middleware"
	"github.com/kholodihor/cows-shelter-backend/models"
)

// CreatePartnerRequest represents the JSON request body for creating a partner
type CreatePartnerRequest struct {
	Name     string `json:"name" binding:"required"`
	Link     string `json:"link"`
	LogoData string `json:"logo_data" binding:"required"` // base64-encoded image data
}

// UpdatePartnerRequest represents the JSON request body for updating a partner
type UpdatePartnerRequest struct {
	Name     string `json:"name"`
	Link     string `json:"link"`
	LogoData string `json:"logo_data"` // base64-encoded image data (optional)
}

func GetAllPartners(c *gin.Context) {
	partners := []models.Partner{}
	config.DB.Find(&partners)
	c.JSON(http.StatusOK, &partners)
}

func GetPartners(c *gin.Context) {
	var partners []models.Partner
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
	config.DB.Model(&models.Partner{}).Count(&total)

	// Fetch the paginated results
	if err := config.DB.Limit(limit).Offset(offset).Find(&partners).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching partners"})
		return
	}

	// Return paginated response
	c.JSON(http.StatusOK, gin.H{
		"data":       partners,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + int64(limit) - 1) / int64(limit), // Calculate total pages
	})
}

func GetPartnerByID(c *gin.Context) {
	var partner models.Partner
	id := c.Param("id")
	if err := config.DB.Where("id = ?", id).First(&partner).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Partner not found"})
		return
	}
	c.JSON(http.StatusOK, &partner)
}

// CreatePartner handles the creation of a new partner with a base64-encoded logo
// @Summary Create a new partner
// @Description Create a new partner with a base64-encoded logo
// @Tags partners
// @Accept json
// @Produce json
// @Param input body CreatePartnerRequest true "Partner data"
// @Success 201 {object} models.Partner
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /partners [post]
func CreatePartner(c *gin.Context) {
	var req CreatePartnerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	// Get storage service from context
	store := middleware.GetStorage(c.Request.Context())
	if store == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage service not available"})
		return
	}

	// Upload the logo using the storage service
	imageURL, err := store.UploadBase64(c.Request.Context(), req.LogoData, "partners")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload logo: " + err.Error()})
		return
	}

	// Create the partner with the logo URL
	partner := models.Partner{
		Name: req.Name,
		Link: req.Link,
		Logo: imageURL,
	}

	// Save the partner to the database
	if err := config.DB.Create(&partner).Error; err != nil {
		// Attempt to delete the uploaded image if database operation fails
		_ = store.DeleteFile(c.Request.Context(), store.ExtractObjectName(imageURL))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create partner"})
		return
	}

	c.JSON(http.StatusCreated, &partner)
}

// UpdatePartner handles updating an existing partner with an optional new base64-encoded logo
// @Summary Update an existing partner
// @Description Update an existing partner with an optional new base64-encoded logo
// @Tags partners
// @Accept json
// @Produce json
// @Param id path int true "Partner ID"
// @Param input body UpdatePartnerRequest true "Updated partner data"
// @Success 200 {object} models.Partner
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /partners/{id} [put]
func UpdatePartner(c *gin.Context) {
	// Get the partner ID from the URL
	id := c.Param("id")
	var partner models.Partner
	if err := config.DB.Where("id = ?", id).First(&partner).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Partner not found"})
		return
	}

	// Parse the request body
	var req UpdatePartnerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	// Update the partner fields if they are provided in the request
	if req.Name != "" {
		partner.Name = req.Name
	}
	if req.Link != "" {
		partner.Link = req.Link
	}

	// Handle logo update if new image data is provided
	if req.LogoData != "" {
		// Get storage service from context
		store := middleware.GetStorage(c.Request.Context())
		if store == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage service not available"})
			return
		}

		// Delete the old logo if it exists
		if partner.Logo != "" {
			oldObjectName := store.ExtractObjectName(partner.Logo)
			if oldObjectName != "" {
				_ = store.DeleteFile(c.Request.Context(), oldObjectName)
			}
		}

		// Upload the new logo using the storage service
		imageURL, err := store.UploadBase64(c.Request.Context(), req.LogoData, "partners")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload new logo: " + err.Error()})
			return
		}

		// Set the new logo URL
		partner.Logo = imageURL
	}

	// Save the updated partner to the database
	if err := config.DB.Save(&partner).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update partner"})
		return
	}

	c.JSON(http.StatusOK, &partner)
}

func DeletePartner(c *gin.Context) {
	// Get the partner ID from the URL
	id := c.Param("id")
	var partner models.Partner
	if err := config.DB.Where("id = ?", id).First(&partner).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Partner not found"})
		return
	}

	// Get storage service from context
	store := middleware.GetStorage(c.Request.Context())
	if store == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage service not available"})
		return
	}

	// Delete the logo if it exists
	if partner.Logo != "" {
		objectName := store.ExtractObjectName(partner.Logo)
		if objectName != "" {
			_ = store.DeleteFile(c.Request.Context(), objectName)
		}
	}

	// Delete the partner from the database
	if err := config.DB.Delete(&partner).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete partner"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Partner deleted successfully"})
}