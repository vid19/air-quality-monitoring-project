package main

import (
	"air-quality/backend/config"
	"air-quality/backend/handler"
	"air-quality/backend/middleware"
	"air-quality/backend/models"

	"github.com/gin-gonic/gin"
)

func main() {
	// config.init() runs automatically, connecting to DB
	
	// This creates the 'users' table from our model
	config.DB.AutoMigrate(&models.User{}, &models.Location{})

	//Set up the Gin router
	r := gin.Default()

	// Group API routes
	api := r.Group("/api")
	{
		// Public Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", handler.RegisterHandler)
			auth.POST("/login", handler.LoginHandler)
		}

		// Public Data routes
		api.GET("/current-weather", handler.WeatherHandler)
		api.GET("/prediction", handler.PredictionHandler)

		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("/locations", handler.AddLocationHandler)
			protected.GET("/locations", handler.GetLocationsHandler)
			protected.DELETE("/locations/:id", handler.DeleteLocationHandler)
		}
	}

	r.Run()
}