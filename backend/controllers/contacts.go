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

	// Find the existing contact
	if err := config.DB.Where("id = ?", id).First(&contact).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
		return
	}

	// Use a map to handle partial updates
	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Update only the fields that are provided in the request
	if name, ok := updateData["name"].(string); ok {
		contact.Name = name
	}

	if email, ok := updateData["email"].(string); ok {
		// Validate email if provided
		if email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email cannot be empty"})
			return
		}
		contact.Email = email
	}

	if phone, ok := updateData["phone"].(string); ok {
		// Validate phone if provided
		if phone == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Phone cannot be empty"})
			return
		}
		contact.Phone = phone
	}

	// Save the updated contact
	if err := config.DB.Save(&contact).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating contact"})
		return
	}

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
