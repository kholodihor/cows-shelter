package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kholodihor/cows-shelter-backend/config"
	"github.com/kholodihor/cows-shelter-backend/controllers"
	"github.com/kholodihor/cows-shelter-backend/middleware"
	"github.com/kholodihor/cows-shelter-backend/storage"
)

type Handler struct{}

type Config struct {
	R       *gin.Engine
	Storage storage.Service
}

// HealthCheckResponse represents the health check response structure
type HealthCheckResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Database  string    `json:"database,omitempty"`
	Storage   string    `json:"storage,omitempty"`
}

// HealthCheck performs health checks for the application
func HealthCheck(c *gin.Context) {
	healthStatus := HealthCheckResponse{
		Status:    "UP",
		Timestamp: time.Now().UTC(),
	}

	// Check database connection
	db := config.DB
	if db == nil {
		healthStatus.Database = "NOT CONFIGURED"
		healthStatus.Status = "DEGRADED"
	} else {
		sqlDB, err := db.DB()
		if err != nil {
			healthStatus.Database = fmt.Sprintf("ERROR: %v", err)
			healthStatus.Status = "DEGRADED"
		} else {
			if err := sqlDB.Ping(); err != nil {
				healthStatus.Database = fmt.Sprintf("DOWN: %v", err)
				healthStatus.Status = "DEGRADED"
			} else {
				healthStatus.Database = "UP"
			}
		}
	}

	// Check storage service
	storageSvc := middleware.GetStorage(c.Request.Context())
	if storageSvc == nil {
		healthStatus.Storage = "NOT CONFIGURED: Storage service not initialized"
		healthStatus.Status = "DEGRADED"
	} else {
		// Simple storage health check - try to list objects
		_, err := storageSvc.ListObjects(c.Request.Context(), "", 1)
		if err != nil {
			healthStatus.Storage = fmt.Sprintf("DOWN: %v", err)
			healthStatus.Status = "DEGRADED"
		} else {
			healthStatus.Storage = "UP"
		}
	}

	// Determine final status
	if healthStatus.Database == "UP" && healthStatus.Storage == "UP" && healthStatus.Status != "DEGRADED" {
		healthStatus.Status = "UP"
	} else if healthStatus.Status != "DEGRADED" {
		healthStatus.Status = "DEGRADED"
	}

	c.JSON(http.StatusOK, healthStatus)
}

func RouterHandler(c *Config) {
	// Health check endpoint
	c.R.GET("/health", HealthCheck)

	// API routes
	c.R.POST("/api/user", controllers.CreateUser)
	c.R.POST("/api/login", controllers.LoginUser)
	c.R.GET("/api/contacts", controllers.GetContacts)
	c.R.GET("/api/excursions/pagination", controllers.GetExcursions)
	c.R.GET("/api/excursions", controllers.GetAllExcursions)
	c.R.GET("/api/excursions/:id", controllers.GetExcursionByID)
	c.R.GET("/api/gallery/pagination", controllers.GetGalleries)
	c.R.GET("/api/gallery", controllers.GetAllGalleries)
	c.R.GET("/api/reviews/pagination", controllers.GetReviews)
	c.R.GET("/api/reviews", controllers.GetAllReviews)
	c.R.GET("/api/reviews/:id", controllers.GetReviewByID)
	c.R.POST("/api/reviews", controllers.CreateReview)
	c.R.PATCH("/api/reviews/:id", controllers.UpdateReview)
	c.R.DELETE("/api/reviews/:id", controllers.DeleteReview)
	c.R.GET("/api/news/pagination", controllers.GetNews)
	c.R.GET("/api/news", controllers.GetAllNews)
	c.R.GET("/api/news/:id", controllers.GetNewsByID)
	c.R.GET("/api/pdf", controllers.GetPdfs)
	c.R.GET("/api/partners", controllers.GetAllPartners)
	c.R.GET("/api/partners/pagination", controllers.GetPartners)
	c.R.GET("/api/partners/:id", controllers.GetPartnerByID)
	c.R.POST("/api/excursions", controllers.CreateExcursion)
	c.R.POST("/api/upload-image", controllers.UploadImage)
	c.R.POST("/api/gallery", controllers.CreateGallery)
	c.R.POST("/api/news", controllers.CreateNews)
	c.R.POST("/api/partners", controllers.CreatePartner)
	c.R.POST("/api/pdf", controllers.CreatePdf)
	c.R.POST("/api/contacts", controllers.CreateContact)
	c.R.PUT("/api/contacts/:id", controllers.UpdateContact)
	c.R.PATCH("/api/contacts/:id", controllers.UpdateContact)
	c.R.PUT("/api/news/:id", controllers.UpdateNews)
	c.R.PUT("/api/excursions/:id", controllers.UpdateExcursion)
	c.R.PATCH("/api/excursions/:id", controllers.UpdateExcursion) // Add PATCH support for frontend compatibility
	c.R.PUT("/api/partners/:id", controllers.UpdatePartner)
	c.R.PATCH("/api/partners/:id", controllers.UpdatePartner) // Add PATCH support for consistency
	c.R.DELETE("/api/news/:id", controllers.DeleteNews)
	c.R.GET("/api/gallery/:id", controllers.GetGalleryByID)
	c.R.DELETE("/api/gallery/:id", controllers.DeleteGallery)
	c.R.DELETE("/api/excursions/:id", controllers.DeleteExcursion)
	c.R.DELETE("/api/partners/:id", controllers.DeletePartner)

	api := c.R.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/user/:id", controllers.GetUserByID)
	}
}
