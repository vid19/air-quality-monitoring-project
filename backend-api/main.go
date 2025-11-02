package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//DB Connection
var db *gorm.DB

//  CONFIGURATION 
var jwtSecretKey = []byte("my_super_secret_key")

// init() runs before main()
func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	//  NEW: CONNECT TO DATABASE 
	dsn := "host=localhost user=postgres password=soupvid dbname=postgres port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Database connection established")

	// This automatically creates the "users" table for us based on the User struct.
	db.AutoMigrate(&User{})
}

//added GORM tags to our User struct
type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
}

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
	var existingUser User
	if err := db.First(&existingUser, "username = ?", req.Username).Error; err == nil {
		// We found a user, so they already exist
		c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create new user
	newUser := User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
	}

	// Save new user to the database
	if err := db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}


func loginHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Check if user exists
	var user User
	if err := db.First(&user, "username = ?", req.Username).Error; err != nil {
		// User not found
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check if password is correct
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		// Passwords don't match
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Create JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}


func weatherHandler(c *gin.Context) {
	apiKey := os.Getenv("OPENWEATHER_API_KEY") 
	lat := "51.5072"; lon := "-0.1276"
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


func main() {
	r := gin.Default()
	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", registerHandler)
			auth.POST("/login", loginHandler)
		}
		api.GET("/current-weather", weatherHandler)
		api.GET("/prediction", predictionHandler)
	}
	r.Run() 
}