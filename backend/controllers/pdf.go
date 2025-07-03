package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kholodihor/cows-shelter-backend/config"
	"github.com/kholodihor/cows-shelter-backend/middleware"
	"github.com/kholodihor/cows-shelter-backend/models"
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

	// Get storage service from context
	store := middleware.GetStorage(c.Request.Context())
	if store == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage service not available"})
		return
	}

	// Upload the document using the storage service
	documentURL, err := store.UploadBase64(c.Request.Context(), req.DocumentData, "documents")
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
		_ = store.DeleteFile(c.Request.Context(), store.ExtractObjectName(documentURL))
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

	// Get storage service from context
	store := middleware.GetStorage(c.Request.Context())
	if store == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage service not available"})
		return
	}

	// Delete the document if it exists
	if pdf.DocumentUrl != "" {
		objectName := store.ExtractObjectName(pdf.DocumentUrl)
		if objectName != "" {
			_ = store.DeleteFile(c.Request.Context(), objectName)
		}
	}

	// Delete the PDF from the database
	config.DB.Delete(&pdf)
	c.JSON(http.StatusOK, gin.H{"message": "PDF deleted successfully"})
}
