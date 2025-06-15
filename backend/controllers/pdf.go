package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kholodihor/cows-shelter-backend/config"
	"github.com/kholodihor/cows-shelter-backend/models"
	"github.com/kholodihor/cows-shelter-backend/services"
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
	// Get the title from form data
	title := c.PostForm("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	// Check if a document file was uploaded
	fileHeader, err := c.FormFile("document")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Document file is required"})
		return
	}

	// Initialize MinIO service
	minioService, err := services.NewMinioService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize MinIO service"})
		return
	}

	// Upload the document to MinIO
	documentURL, err := minioService.UploadFile(
		c.Request.Context(),
		fileHeader,
		"documents",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload document: " + err.Error()})
		return
	}

	// Create the new PDF entry
	pdf := models.Pdf{
		Title:       title,
		DocumentUrl: documentURL,
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

	// Initialize MinIO service
	minioService, err := services.NewMinioService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize MinIO service"})
		return
	}

	// Delete the document from MinIO if it exists
	if pdf.DocumentUrl != "" {
		objectName := minioService.ExtractObjectName(pdf.DocumentUrl)
		if objectName != "" {
			_ = minioService.DeleteFile(c.Request.Context(), objectName)
		}
	}

	// Delete the PDF from the database
	config.DB.Delete(&pdf)
	c.JSON(http.StatusOK, gin.H{"message": "PDF deleted successfully"})
}
