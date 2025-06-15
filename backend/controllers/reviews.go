package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kholodihor/cows-shelter-backend/config"
	"github.com/kholodihor/cows-shelter-backend/models"
)

func GetAllReviews(c *gin.Context) {
	reviews := []models.Review{}
	config.DB.Find(&reviews)
	c.JSON(http.StatusOK, &reviews)
}

// GetReviews - Retrieve all reviews
func GetReviews(c *gin.Context) {
	var reviews []models.Review
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
	config.DB.Model(&models.Review{}).Count(&total)

	// Fetch the paginated results
	if err := config.DB.Limit(limit).Offset(offset).Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching reviews"})
		return
	}

	// Return paginated response
	c.JSON(http.StatusOK, gin.H{
		"data":       reviews,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + int64(limit) - 1) / int64(limit), // Calculate total pages
	})
}


// GetReviewByID - Retrieve a specific review by ID
func GetReviewByID(c *gin.Context) {
	var review models.Review
	id := c.Param("id")

	if err := config.DB.Where("id = ?", id).First(&review).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}
	c.JSON(http.StatusOK, &review)
}

// CreateReview - Create a new review entry
func CreateReview(c *gin.Context) {
	var requestBody struct {
		NameEn   string `json:"name_en" binding:"required"`
		NameUa   string `json:"name_ua" binding:"required"`
		ReviewEn string `json:"review_en" binding:"required"`
		ReviewUa string `json:"review_ua" binding:"required"`
	}

	// Validate the request body
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Create the new review entry
	review := models.Review{
		NameEn:   requestBody.NameEn,
		NameUa:   requestBody.NameUa,
		ReviewEn: requestBody.ReviewEn,
		ReviewUa: requestBody.ReviewUa,
		CreatedAt: time.Now(),
	}

	if err := config.DB.Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating review"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         review.ID,
		"name_en":    review.NameEn,
		"name_ua":    review.NameUa,
		"review_en":  review.ReviewEn,
		"review_ua":  review.ReviewUa,
		"created_at": review.CreatedAt,
	})
}

// UpdateReview - Update a specific review by ID
func UpdateReview(c *gin.Context) {
	id := c.Param("id")
	var review models.Review

	// Check if review exists
	if err := config.DB.Where("id = ?", id).First(&review).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	// Define a struct to hold the updateable fields
	var updateData struct {
		NameEn   *string `json:"name_en,omitempty"`
		NameUa   *string `json:"name_ua,omitempty"`
		ReviewEn *string `json:"review_en,omitempty"`
		ReviewUa *string `json:"review_ua,omitempty"`
	}

	// Bind the request body to the updateData struct
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Update only the provided fields
	updates := make(map[string]interface{})
	if updateData.NameEn != nil {
		updates["name_en"] = *updateData.NameEn
	}
	if updateData.NameUa != nil {
		updates["name_ua"] = *updateData.NameUa
	}
	if updateData.ReviewEn != nil {
		updates["review_en"] = *updateData.ReviewEn
	}
	if updateData.ReviewUa != nil {
		updates["review_ua"] = *updateData.ReviewUa
	}

	// Update the review in the database
	if err := config.DB.Model(&review).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update review"})
		return
	}

	// Fetch the updated review
	var updatedReview models.Review
	if err := config.DB.First(&updatedReview, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated review"})
		return
	}

	c.JSON(http.StatusOK, updatedReview)
}

// DeleteReview - Delete a specific review by ID
func DeleteReview(c *gin.Context) {
	id := c.Param("id")
	var review models.Review

	// Check if review exists before deleting
	if err := config.DB.Where("id = ?", id).First(&review).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	// Delete the review
	if err := config.DB.Delete(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}