package main

import (
	"air-quality/backend/config"
	"air-quality/backend/handler"
	"air-quality/backend/models"

	"github.com/gin-gonic/gin"
)

func main() {
	// config.init() runs automatically, connecting to DB

	// Run auto-migration
	// This creates the 'users' table from our model
	config.DB.AutoMigrate(&models.User{})

	// Set up the Gin router
	r := gin.Default()

	// Group API routes
	api := r.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", handler.RegisterHandler)
			auth.POST("/login", handler.LoginHandler)
		}

		// Data routes
		api.GET("/current-weather", handler.WeatherHandler)
		api.GET("/prediction", handler.PredictionHandler)
	}

	r.Run()
}