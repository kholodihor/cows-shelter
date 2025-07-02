package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kholodihor/cows-shelter-backend/config"
	"github.com/kholodihor/cows-shelter-backend/models"
	"github.com/kholodihor/cows-shelter-backend/services"
)

func GetAllNews(c *gin.Context) {
	news := []models.News{}
	config.DB.Find(&news)
	c.JSON(http.StatusOK, &news)
}

func GetNews(c *gin.Context) {
	var news []models.News
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
	config.DB.Model(&models.News{}).Count(&total)

	// Fetch the paginated results
	if err := config.DB.Limit(limit).Offset(offset).Find(&news).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching news"})
		return
	}

	// Return paginated response
	c.JSON(http.StatusOK, gin.H{
		"data":       news,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + int64(limit) - 1) / int64(limit), // Calculate total pages
	})
}

func GetNewsByID(c *gin.Context) {
	var news models.News
	id := c.Param("id")
	if err := config.DB.Where("id = ?", id).First(&news).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
		return
	}
	c.JSON(http.StatusOK, &news)
}

// CreateNewsRequest represents the JSON request body for creating a news item
type CreateNewsRequest struct {
	TitleEn    string `json:"title_en" binding:"required"`
	TitleUa    string `json:"title_ua"`
	SubtitleEn string `json:"subtitle_en"`
	SubtitleUa string `json:"subtitle_ua"`
	ContentEn  string `json:"content_en" binding:"required"`
	ContentUa  string `json:"content_ua"`
	ImageData  string `json:"image_data"` // base64-encoded image data
}

// CreateNews handles the creation of a news item with an optional image
// @Summary Create a new news item
// @Description Create a new news item with an optional base64-encoded image
// @Tags news
// @Accept json
// @Produce json
// @Param input body CreateNewsRequest true "News data"
// @Success 201 {object} models.News
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /news [post]
func CreateNews(c *gin.Context) {
	var req CreateNewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	// Create news model
	news := models.News{
		TitleEn:    req.TitleEn,
		TitleUa:    req.TitleUa,
		SubtitleEn: req.SubtitleEn,
		SubtitleUa: req.SubtitleUa,
		ContentEn:  req.ContentEn,
		ContentUa:  req.ContentUa,
	}

	// Handle image upload if present
	if req.ImageData != "" {
		// Initialize MinIO service
		minioService, err := services.NewMinioService()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize MinIO service"})
			return
		}

		// Upload the base64 image to MinIO
		imageURL, err := minioService.UploadBase64(c.Request.Context(), req.ImageData, "news")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image: " + err.Error()})
			return
		}
		news.ImageUrl = imageURL
	}

	// Save the news entry to the database
	if err := config.DB.Create(&news).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create news entry"})
		return
	}

	c.JSON(http.StatusCreated, news)
}


// UpdateNewsRequest represents the JSON request body for updating a news item
type UpdateNewsRequest struct {
	TitleEn    string `json:"title_en"`
	TitleUa    string `json:"title_ua"`
	SubtitleEn string `json:"subtitle_en"`
	SubtitleUa string `json:"subtitle_ua"`
	ContentEn  string `json:"content_en"`
	ContentUa  string `json:"content_ua"`
	ImageData  string `json:"image_data"` // base64-encoded image data
}

// UpdateNews handles updating a news item with an optional new image
// @Summary Update a news item
// @Description Update a news item with an optional new base64-encoded image
// @Tags news
// @Accept json
// @Produce json
// @Param id path int true "News ID"
// @Param input body UpdateNewsRequest true "Updated news data"
// @Success 200 {object} models.News
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /news/{id} [put]
func UpdateNews(c *gin.Context) {
	// Get the existing news item
	var news models.News
	if err := config.DB.Where("id = ?", c.Param("id")).First(&news).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
		return
	}

	var req UpdateNewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	// Update fields if provided
	if req.TitleEn != "" {
		news.TitleEn = req.TitleEn
	}
	if req.TitleUa != "" {
		news.TitleUa = req.TitleUa
	}
	if req.SubtitleEn != "" {
		news.SubtitleEn = req.SubtitleEn
	}
	if req.SubtitleUa != "" {
		news.SubtitleUa = req.SubtitleUa
	}
	if req.ContentEn != "" {
		news.ContentEn = req.ContentEn
	}
	if req.ContentUa != "" {
		news.ContentUa = req.ContentUa
	}

	// Handle image upload if new image data is provided
	if req.ImageData != "" {
		// Initialize MinIO service
		minioService, err := services.NewMinioService()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize MinIO service"})
			return
		}

		// Upload the new base64 image to MinIO
		newImageURL, err := minioService.UploadBase64(c.Request.Context(), req.ImageData, "news")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload new image: " + err.Error()})
			return
		}

		// Delete the old image from MinIO if it exists
		if news.ImageUrl != "" {
			oldObjectName := minioService.ExtractObjectName(news.ImageUrl)
			if oldObjectName != "" {
				_ = minioService.DeleteFile(c.Request.Context(), oldObjectName)
			}
		}

		// Update the image URL
		news.ImageUrl = newImageURL
	}

	// Save the updated news item
	if err := config.DB.Save(&news).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update news"})
		return
	}

	c.JSON(http.StatusOK, &news)
}

// DeleteNews handles the deletion of a news item and its associated image
// @Summary Delete a news item
// @Description Delete a news item and its associated image
// @Tags news
// @Produce json
// @Param id path int true "News ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /news/{id} [delete]
func DeleteNews(c *gin.Context) {
	// Get the news item to delete
	var news models.News
	if err := config.DB.Where("id = ?", c.Param("id")).First(&news).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
		return
	}

	// Initialize MinIO service
	minioService, err := services.NewMinioService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize MinIO service"})
		return
	}

	// Delete the image from MinIO if it exists
	if news.ImageUrl != "" {
		// Extract the object name from the URL
		objectName := minioService.ExtractObjectName(news.ImageUrl)
		if objectName != "" {
			_ = minioService.DeleteFile(c.Request.Context(), objectName)
		}
	}

	// Delete the news item from the database
	if err := config.DB.Delete(&news).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete news"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "News item deleted successfully"})
}