package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kholodihor/cows-shelter-backend/config"
	"github.com/kholodihor/cows-shelter-backend/models"
)

// GetContacts - Retrieve all contacts
func GetContacts(c *gin.Context) {
	var contacts []models.Contact
	config.DB.Find(&contacts)
	c.JSON(http.StatusOK, &contacts)
}

// GetContactByID - Retrieve a specific contact by ID
func GetContactByID(c *gin.Context) {
	var contact models.Contact
	id := c.Param("id")

	if err := config.DB.Where("id = ?", id).First(&contact).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
		return
	}
	c.JSON(http.StatusOK, &contact)
}

// CreateContact - Create a new contact entry
// @Summary Create a new contact
// @Description Create a new contact with the provided details
// @Tags contacts
// @Accept json
// @Produce json
// @Param contact body models.Contact true "Contact details"
// @Success 201 {object} models.Contact
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /contacts [post]
func CreateContact(c *gin.Context) {
	var contact models.Contact

	// Bind JSON to contact model
	if err := c.ShouldBindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Validate required fields
	if contact.Name == "" || contact.Email == "" || contact.Phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name, email, and phone are required"})
		return
	}

	// Save to database
	if err := config.DB.Create(&contact).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating contact"})
		return
	}

	c.JSON(http.StatusCreated, contact)
}

// UpdateContact - Update an existing contact
func UpdateContact(c *gin.Context) {
	id := c.Param("id")
	var contact models.Contact

	if err := config.DB.Where("id = ?", id).First(&contact).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
		return
	}

	var requestBody struct {
		Email string `json:"email" binding:"required,email"`
		Phone string `json:"phone" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	contact.Email = requestBody.Email
	contact.Phone = requestBody.Phone

	config.DB.Save(&contact)
	c.JSON(http.StatusOK, &contact)
}

// DeleteContact - Delete a specific contact by ID
func DeleteContact(c *gin.Context) {
	id := c.Param("id")
	var contact models.Contact

	// Check if contact exists before deleting
	if err := config.DB.Where("id = ?", id).First(&contact).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
		return
	}

	// Delete the contact
	config.DB.Delete(&contact)
	c.JSON(http.StatusOK, gin.H{"message": "Contact deleted successfully"})
}
