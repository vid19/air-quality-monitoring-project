package main

import (
	"io/ioutil"
	"net/http"
	"time"

	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)
func init() {
	// Load the .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}
}

// --- MOCK USER DATABASE ---
// In a real app, this would be your PostgreSQL database
var userDatabase = make(map[string]User)

// This is our "User" model
type User struct {
	Username     string `json:"username"`
	PasswordHash string `json:"-"` // Don't send the hash in JSON
}

// --- CONFIGURATION ---
// IMPORTANT: Change this to a long, random secret key!
var jwtSecretKey = []byte("my_super_secret_key")

// --- AUTH HANDLERS ---

// RegisterRequest holds the data from a user's registration
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func registerHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Check if user already exists
	if _, exists := userDatabase[req.Username]; exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Save new user to our mock database
	userDatabase[req.Username] = User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func loginHandler(c *gin.Context) {
	var req RegisterRequest // Re-use the same struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Check if user exists
	user, exists := userDatabase[req.Username]
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check if password is correct
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		// Passwords don't match
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// --- Create a JWT Token ---
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})

	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	// Send the token back
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func weatherHandler(c *gin.Context) {
	apiKey := os.Getenv("OPENWEATHER_API_KEY") // 👈 CHANGED
	if apiKey == "" {
		log.Println("API key not set")
	}

	lat := "51.5072"
	lon := "-0.1276"
	url := "https://api.openweathermap.org/data/2.5/air_pollution?lat=" + lat + "&lon=" + lon + "&appid=" + apiKey
	resp, err := http.Get(url)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data from OpenWeather"}); return }
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"}); return }
	c.Data(http.StatusOK, "application/json", body)
}

func predictionHandler(c *gin.Context) {
	pythonServiceUrl := "http://localhost:8001/predict"
	resp, err := http.Post(pythonServiceUrl, "application/json", nil)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch prediction from AI service"}); return }
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read AI response body"}); return }
	c.Data(http.StatusOK, "application/json", body)
}


// --- MAIN FUNCTION ---
func main() {
	r := gin.Default()

	// Group API routes
	api := r.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", registerHandler)
			auth.POST("/login", loginHandler)
		}

		// Data routes
		api.GET("/current-weather", weatherHandler)
		api.GET("/prediction", predictionHandler)
		
		// We'll add a "protected" route here later
	}

	r.Run() // Runs on http://localhost:8080
}