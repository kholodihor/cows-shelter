package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kholodihor/cows-shelter-backend/config"
	"github.com/kholodihor/cows-shelter-backend/models"
	"github.com/kholodihor/cows-shelter-backend/utils"
)

func LoginUser(c *gin.Context) {
	var requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", requestBody.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email"})
		return
	}

	// Check if password is correct
	if !utils.CheckPasswordHash(requestBody.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	// Return token and email
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"email": user.Email,
	})
}