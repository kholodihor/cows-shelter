package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kholodihor/cows-shelter-backend/config"
	"github.com/kholodihor/cows-shelter-backend/models"
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
	var partner models.Partner

	// Bind JSON to the partner model, which includes other details
	if err := c.ShouldBindJSON(&partner); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
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
		"logo": partner.Logo, // This will be set after the image is uploaded
	})
}

func UpdatePartner(c *gin.Context) {
	var partner models.Partner
	if err := config.DB.Where("id = ?", c.Param("id")).First(&partner).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Partner not found"})
		return
	}

	// Parse the file if a new logo is uploaded
	file, err := c.FormFile("logo")
	if err == nil {
		// Save the new file
		uploadDir := "./uploads/partners"
		filename := time.Now().Format("20060102150405") + "_" + filepath.Base(file.Filename)
		filepath := filepath.Join(uploadDir, filename)
		if err := c.SaveUploadedFile(file, filepath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save file"})
			return
		}

		// Optionally, delete the old logo file
		if err := os.Remove("." + partner.Logo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete old logo"})
			return
		}

		// Set the new logo URL
		partner.Logo = "/uploads/partners/" + filename
	}

	if err := c.ShouldBind(&partner); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
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

	if err := config.DB.Delete(&partner).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete partner"})
		return
	}

	// Optionally, delete the logo file from the server
	if err := os.Remove("." + partner.Logo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete logo file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Partner deleted"})
}