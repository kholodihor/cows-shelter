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

// CreateExcursion handles the creation of a new excursion with an optional image
// @Summary Create a new excursion
// @Description Create a new excursion with an optional image
// @Tags excursions
// @Accept multipart/form-data
// @Produce json
// @Param title_en formData string true "Title in English"
// @Param title_ua formData string false "Title in Ukrainian"
// @Param description_en formData string true "Description in English"
// @Param description_ua formData string false "Description in Ukrainian"
// @Param time_to formData string true "End time of the excursion"
// @Param time_from formData string true "Start time of the excursion"
// @Param amount_of_persons formData string true "Maximum number of persons"
// @Param image formData file false "Excursion image"
// @Success 201 {object} models.Excursion
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /excursions [post]
func CreateExcursion(c *gin.Context) {
	// Parse the multipart form
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10 MB max size
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form: " + err.Error()})
		return
	}

	// Get form values
	excursion := models.Excursion{
		TitleEn:         c.PostForm("title_en"),
		TitleUa:         c.PostForm("title_ua"),
		DescriptionEn:   c.PostForm("description_en"),
		DescriptionUa:   c.PostForm("description_ua"),
		TimeTo:          c.PostForm("time_to"),
		TimeFrom:        c.PostForm("time_from"),
		AmountOfPersons: c.PostForm("amount_of_persons"),
		CreatedAt:       time.Now(),
	}

	// Validate required fields
	if excursion.TitleEn == "" || excursion.DescriptionEn == "" || excursion.TimeTo == "" || excursion.TimeFrom == "" || excursion.AmountOfPersons == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title, description, time_to, time_from, and amount_of_persons are required"})
		return
	}

	// Handle image upload if present
	file, fileHeader, err := c.Request.FormFile("image")
	if err == nil {
		defer file.Close()

		// Initialize MinIO service
		minioService, err := services.NewMinioService()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize MinIO service"})
			return
		}

		// Upload the image to MinIO
		imageURL, err := minioService.UploadFile(
			c.Request.Context(),
			fileHeader,
			"excursions",
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image: " + err.Error()})
			return
		}
		excursion.ImageUrl = imageURL
	}

	// Save the excursion to the database
	if err := config.DB.Create(&excursion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create excursion"})
		return
	}

	c.JSON(http.StatusCreated, excursion)
}

// UpdateExcursion handles updating an existing excursion with an optional new image
// @Summary Update an existing excursion
// @Description Update an existing excursion with an optional new image
// @Tags excursions
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Excursion ID"
// @Param title_en formData string false "Title in English"
// @Param title_ua formData string false "Title in Ukrainian"
// @Param description_en formData string false "Description in English"
// @Param description_ua formData string false "Description in Ukrainian"
// @Param time_to formData string false "End time of the excursion"
// @Param time_from formData string false "Start time of the excursion"
// @Param amount_of_persons formData string false "Maximum number of persons"
// @Param image formData file false "New excursion image"
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

	// Parse the multipart form
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10 MB max size
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form: " + err.Error()})
		return
	}

	// Update fields if provided
	if titleEn := c.PostForm("title_en"); titleEn != "" {
		excursion.TitleEn = titleEn
	}
	if titleUa := c.PostForm("title_ua"); titleUa != "" {
		excursion.TitleUa = titleUa
	}
	if descriptionEn := c.PostForm("description_en"); descriptionEn != "" {
		excursion.DescriptionEn = descriptionEn
	}
	if descriptionUa := c.PostForm("description_ua"); descriptionUa != "" {
		excursion.DescriptionUa = descriptionUa
	}
	if timeTo := c.PostForm("time_to"); timeTo != "" {
		excursion.TimeTo = timeTo
	}
	if timeFrom := c.PostForm("time_from"); timeFrom != "" {
		excursion.TimeFrom = timeFrom
	}
	if amountOfPersons := c.PostForm("amount_of_persons"); amountOfPersons != "" {
		excursion.AmountOfPersons = amountOfPersons
	}

	// Initialize MinIO service
	minioService, err := services.NewMinioService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize MinIO service"})
		return
	}

	// Handle image upload if present
	file, fileHeader, err := c.Request.FormFile("image")
	if err == nil {
		defer file.Close()

		// Upload the new image to MinIO
		newImageURL, err := minioService.UploadFile(
			c.Request.Context(),
			fileHeader,
			"excursions",
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload new image: " + err.Error()})
			return
		}

		// Delete the old image from MinIO if it exists
		if excursion.ImageUrl != "" {
			oldObjectName := minioService.ExtractObjectName(excursion.ImageUrl)
			if oldObjectName != "" {
				_ = minioService.DeleteFile(c.Request.Context(), oldObjectName)
			}
		}

		// Update the image URL
		excursion.ImageUrl = newImageURL
	}

	// Save the updated excursion
	if err := config.DB.Save(&excursion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update excursion"})
		return
	}

	c.JSON(http.StatusOK, &excursion)
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

	// Initialize MinIO service
	minioService, err := services.NewMinioService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize MinIO service"})
		return
	}

	// Delete the image from MinIO if it exists
	if excursion.ImageUrl != "" {
		// Extract the object name from the URL
		objectName := minioService.ExtractObjectName(excursion.ImageUrl)
		if objectName != "" {
			_ = minioService.DeleteFile(c.Request.Context(), objectName)
		}
	}

	// Delete the excursion from the database
	if err := config.DB.Delete(&excursion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete excursion"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Excursion deleted successfully"})
}