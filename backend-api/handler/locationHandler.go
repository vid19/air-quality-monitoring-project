package handler

import (
	"air-quality/backend/config"
	"air-quality/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddLocationHandler(c *gin.Context) {
	var req models.AddLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	userID, _ := c.Get("userID")

	// Create the new location model
	newLocation := models.Location{
		Name:      req.Name,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		UserID:    userID.(uint),
	}

	// Save to database
	if err := config.DB.Create(&newLocation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save location"})
		return
	}

	c.JSON(http.StatusCreated, newLocation)
}

// GetLocationsHandler retrieves all locations for the logged-in user
func GetLocationsHandler(c *gin.Context) {
	// Get the user ID from the context
	userID, _ := c.Get("userID")

	var locations []models.Location

	// Find all locations in the DB where user_id matches
	if err := config.DB.Preload("User").Where("user_id = ?", userID).Find(&locations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve locations"})
		return
	}

	c.JSON(http.StatusOK, locations)
}

// DeleteLocationHandler deletes a location for the logged-in user
func DeleteLocationHandler(c *gin.Context) {
    // Get the location ID from the URL (e.g., /api/locations/5)
    locationID := c.Param("id")

    // Get the user ID from the context
    userID, _ := c.Get("userID")

    var location models.Location
    // First, find the location
    if err := config.DB.First(&location, locationID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Location not found"})
        return
    }

    // CRITICAL: Check if this user actually owns this location
    if location.UserID != userID.(uint) {
        c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this location"})
        return
    }

    // User owns it, so delete it
    if err := config.DB.Delete(&location).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete location"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Location deleted successfully"})
}