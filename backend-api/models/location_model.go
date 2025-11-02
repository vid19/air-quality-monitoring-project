package models

type Location struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	Name      string  `gorm:"not null" json:"name"`
	Latitude  float64 `gorm:"not null" json:"latitude"`
	Longitude float64 `gorm:"not null" json:"longitude"`
	UserID    uint    `gorm:"not null" json:"-"` //foreign key
	User      User    `json:"user"`
}

// This struct we'll expect from the frontend
type AddLocationRequest struct {
	Name      string  `json:"name" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}