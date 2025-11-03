package handler

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// getLatLon gets latitude and longitude from the query or returns a default (London)
func getLatLon(c *gin.Context) (string, string) {
    lat := c.DefaultQuery("lat", "51.5072") // Default to London
    lon := c.DefaultQuery("lon", "-0.1276") // Default to London
    return lat, lon
}

func WeatherHandler(c *gin.Context) {
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		log.Println("API key not set")
	}
    
    //Get DYNAMIC LAT/LON
    lat, lon := getLatLon(c)
	
	url := "https://api.openweathermap.org/data/2.5/air_pollution?lat=" + lat + "&lon=" + lon + "&appid=" + apiKey
	
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data from OpenWeather"})
		return
	}
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}
	c.Data(http.StatusOK, "application/json", body)
}

func PredictionHandler(c *gin.Context) {
    //GET DYNAMIC LAT/LON
    // Even though our "mock" AI doesn't use this yet,
    // the Go handler should still support it for when we build the real AI.
    // For now, it just demonstrates the data flow.
    lat, lon := getLatLon(c)
    log.Printf("AI Prediction requested for lat: %s, lon: %s", lat, lon)

	pythonServiceUrl := "http://localhost:8001/predict"
	
    //need to update Python service later to actually use this data.For now, we make a simple POST.
	resp, err := http.Post(pythonServiceUrl, "application/json", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch prediction from AI service"})
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read AI response body"})
		return
	}
	c.Data(http.StatusOK, "application/json", body)
}