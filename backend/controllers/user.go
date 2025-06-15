package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kholodihor/cows-shelter-backend/config"
	"github.com/kholodihor/cows-shelter-backend/models"
	"github.com/kholodihor/cows-shelter-backend/utils"
	"gorm.io/gorm"
)

func GetUsers(c *gin.Context){
	users := []models.User{}
	config.DB.Find(&users)
	c.JSON(http.StatusOK,&users)
}

func GetUserByID(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	if err := config.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, &user)
}

func CreateUser(c *gin.Context) {
    var requestBody struct {
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required,min=8"`
    }

    if err := c.ShouldBindJSON(&requestBody); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
        return
    }

    // Check if the email is already taken
    var existingUser models.User
    if err := config.DB.Where("email = ?", requestBody.Email).First(&existingUser).Error; err == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Email already taken"})
        return
    } else if !errors.Is(err, gorm.ErrRecordNotFound) {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking existing user"})
        return
    }

    // Hash the password before saving
    hashedPassword, err := utils.HashPassword(requestBody.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
        return
    }

    user := models.User{
        Email:    requestBody.Email,
        Password: hashedPassword,
    }

    // Use a transaction to ensure atomicity
    tx := config.DB.Begin()
    if err := tx.Create(&user).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
        return
    }

    // Generate JWT token
    token, err := utils.GenerateJWT(user.Email)
    if err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
        return
    }

    tx.Commit()

    // Return the token and user ID (not the full user object)
    c.JSON(http.StatusCreated, gin.H{
        "id":    user.ID,
        "email": user.Email,
        "token": token,
    })
}

func UpdateUser(c *gin.Context){
	var user models.User
	config.DB.Where("id = ?", c.Param("id")).First(&user)
	c.BindJSON(&user)
	config.DB.Save(&user)
	c.JSON(200,&user)
}

func DeleteUser(c *gin.Context){
	var user models.User
	config.DB.Where("id = ?", c.Param("id")).Delete(&user)
	c.JSON(200,&user)
}