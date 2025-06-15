package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kholodihor/cows-shelter-backend/config"
	"github.com/kholodihor/cows-shelter-backend/models"
)

// CreatePassword - Create a new password reset token
func CreatePassword(c *gin.Context) {
	var requestBody struct {
		Email string `json:"email" binding:"required,email"`
	}

	// Validate the request body
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email: " + err.Error()})
		return
	}

	// Generate a unique token
	token := uuid.New().String()

	// Create the new Password entry
	passwordReset := models.Password{
		Email: requestBody.Email,
		Token: token,
	}

	if err := config.DB.Create(&passwordReset).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating password reset token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"email": passwordReset.Email,
		"token": passwordReset.Token,
		"created_at": time.Now(),
	})
}

// GetPasswordByToken - Retrieve a password reset entry by token
func GetPasswordByToken(c *gin.Context) {
	token := c.Param("token")

	var passwordReset models.Password
	if err := config.DB.Where("token = ?", token).First(&passwordReset).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Password reset token not found"})
		return
	}

	c.JSON(http.StatusOK, &passwordReset)
}
