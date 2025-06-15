package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kholodihor/cows-shelter-backend/controllers"
	"github.com/kholodihor/cows-shelter-backend/middleware"
)

type Handler struct{}

type Config struct{
	R *gin.Engine
}

func RouterHandler(c *Config){
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
	c.R.PATCH("/api/excursions/:id", controllers.UpdateExcursion)  // Add PATCH support for frontend compatibility
	c.R.PUT("/api/partners/:id", controllers.UpdatePartner)
	c.R.PATCH("/api/partners/:id", controllers.UpdatePartner)  // Add PATCH support for consistency
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
