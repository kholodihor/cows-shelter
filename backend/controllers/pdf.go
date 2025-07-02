package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kholodihor/cows-shelter-backend/config"
	"github.com/kholodihor/cows-shelter-backend/models"
	"github.com/kholodihor/cows-shelter-backend/services"
)

// CreatePdfRequest represents the JSON request body for creating a PDF
type CreatePdfRequest struct {
	Title      string `json:"title" binding:"required"`
	DocumentData string `json:"document_data" binding:"required"` // base64-encoded document data
}

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

// CreatePdf handles the creation of a new PDF entry with a base64-encoded document
// @Summary Create a new PDF entry
// @Description Create a new PDF entry with a base64-encoded document
// @Tags pdfs
// @Accept json
// @Produce json
// @Param input body CreatePdfRequest true "PDF data"
// @Success 201 {object} models.Pdf
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /pdfs [post]
func CreatePdf(c *gin.Context) {
	// Parse the request body
	var req CreatePdfRequest
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

	// Upload the document to MinIO using base64 data
	documentURL, err := minioService.UploadBase64(c.Request.Context(), req.DocumentData, "documents")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload document: " + err.Error()})
		return
	}

	// Create the new PDF entry
	pdf := models.Pdf{
		Title:       req.Title,
		DocumentUrl: documentURL,
	}

	if err := config.DB.Create(&pdf).Error; err != nil {
		// Attempt to delete the uploaded document if database operation fails
		_ = minioService.DeleteFile(c.Request.Context(), minioService.ExtractObjectName(documentURL))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create PDF"})
		return
	}

	c.JSON(http.StatusCreated, &pdf)
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
