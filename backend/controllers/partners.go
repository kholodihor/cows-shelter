package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kholodihor/cows-shelter-backend/config"
	"github.com/kholodihor/cows-shelter-backend/models"
	"github.com/kholodihor/cows-shelter-backend/services"
)

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

func CreatePartner(c *gin.Context) {
	// Parse form data
	name := c.PostForm("name")
	link := c.PostForm("link")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	// Initialize partner with form data
	partner := models.Partner{
		Name: name,
		Link: link,
	}

	// Check if a logo file was uploaded
	fileHeader, err := c.FormFile("logo")
	if err == nil {
		// Initialize MinIO service
		minioService, err := services.NewMinioService()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize MinIO service"})
			return
		}

		// Upload the logo to MinIO
		imageURL, err := minioService.UploadFile(
			c.Request.Context(),
			fileHeader,
			"partners",
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload logo: " + err.Error()})
			return
		}

		// Set the logo URL in the partner model
		partner.Logo = imageURL
	}

	// Save the partner entry to the database
	if err := config.DB.Create(&partner).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create partner"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":   partner.ID,
		"name": partner.Name,
		"link": partner.Link,
		"logo": partner.Logo,
	})
}

func UpdatePartner(c *gin.Context) {
	var partner models.Partner
	if err := config.DB.Where("id = ?", c.Param("id")).First(&partner).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Partner not found"})
		return
	}

	// Update fields from form data
	if name := c.PostForm("name"); name != "" {
		partner.Name = name
	}

	if link := c.PostForm("link"); link != "" {
		partner.Link = link
	}

	// Parse the file if a new logo is uploaded
	fileHeader, err := c.FormFile("logo")
	if err == nil {
		// Initialize MinIO service
		minioService, err := services.NewMinioService()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize MinIO service"})
			return
		}

		// Delete the old logo from MinIO if it exists
		if partner.Logo != "" {
			oldObjectName := minioService.ExtractObjectName(partner.Logo)
			if oldObjectName != "" {
				_ = minioService.DeleteFile(c.Request.Context(), oldObjectName)
			}
		}

		// Upload the new logo to MinIO
		imageURL, err := minioService.UploadFile(
			c.Request.Context(),
			fileHeader,
			"partners",
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload new logo: " + err.Error()})
			return
		}

		// Set the new logo URL
		partner.Logo = imageURL
	}

	if err := config.DB.Save(&partner).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update partner"})
		return
	}

	c.JSON(http.StatusOK, &partner)
}

func DeletePartner(c *gin.Context) {
	var partner models.Partner
	if err := config.DB.Where("id = ?", c.Param("id")).First(&partner).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Partner not found"})
		return
	}

	// Initialize MinIO service
	minioService, err := services.NewMinioService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize MinIO service"})
		return
	}

	// Delete the logo from MinIO if it exists
	if partner.Logo != "" {
		objectName := minioService.ExtractObjectName(partner.Logo)
		if objectName != "" {
			_ = minioService.DeleteFile(c.Request.Context(), objectName)
		}
	}

	// Delete the partner from the database
	if err := config.DB.Delete(&partner).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete partner"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Partner deleted"})
}