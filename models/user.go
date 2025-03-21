package models

import (
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// User represents a user of the application.
type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Password string
}

// Claims defines the JWT claims structure.
type Claims struct {
	Id        uint   `json:"id"`
	UserEmail string `json:"useremail"`
	jwt.RegisteredClaims
}
