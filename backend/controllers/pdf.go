package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kholodihor/cows-shelter-backend/config"
	"github.com/kholodihor/cows-shelter-backend/models"
)

// GetPdfs - Retrieve all PDFs
func GetPdfs(c *gin.Context) {
	var pdfs []models.Pdf
	config.DB.Find(&pdfs)
	c.JSON(http.StatusOK, &pdfs)
}

// GetPdfByID - Retrieve a specific PDF by ID
func GetPdfByID(c *gin.Context) {
	var pdf models.Pdf
	id := c.Param("id")

	if err := config.DB.Where("id = ?", id).First(&pdf).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "PDF not found"})
		return
	}
	c.JSON(http.StatusOK, &pdf)
}

// CreatePdf - Create a new PDF entry
func CreatePdf(c *gin.Context) {
	var requestBody struct {
		Title       string `json:"title" binding:"required"`
		DocumentUrl string `json:"document_url" binding:"required,url"`
	}

	// Validate the request body
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Create the new PDF entry
	pdf := models.Pdf{
		Title:       requestBody.Title,
		DocumentUrl: requestBody.DocumentUrl,
		CreatedAt:   time.Now(),
	}

	if err := config.DB.Create(&pdf).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating PDF"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":          pdf.ID,
		"title":       pdf.Title,
		"document_url": pdf.DocumentUrl,
		"created_at":  pdf.CreatedAt,
	})
}

// DeletePdf - Delete a specific PDF by ID
func DeletePdf(c *gin.Context) {
	id := c.Param("id")
	var pdf models.Pdf

	// Check if PDF exists before deleting
	if err := config.DB.Where("id = ?", id).First(&pdf).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "PDF not found"})
		return
	}

	// Delete the PDF
	config.DB.Delete(&pdf)
	c.JSON(http.StatusOK, gin.H{"message": "PDF deleted successfully"})
}
