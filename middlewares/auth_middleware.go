package middlewares

import (
	"log"

	"github.com/Govind-619/GoAdminHub/utils"

	"github.com/Govind-619/GoAdminHub/models"

	"net/http"
	"os"

	"github.com/Govind-619/GoAdminHub/config"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// ClearCache sets headers to prevent caching.
func ClearCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Next()
	}
}

// Authenticate middleware for normal users.
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("jwtToken")
		if err != nil || token == "" {
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		email, err := utils.ParseToken(token)
		if err != nil || email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Retrieve user from the database
		var user models.User
		if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("email", email)
		c.Set("user", user)
		c.Next()
	}
}

// AdminAuthenticate middleware for the admin panel.
func AdminAuthenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("Admin")
		if err != nil || token == "" {
			c.Redirect(http.StatusSeeOther, "/admin/login")
			c.Abort()
			return
		}

		email, err := utils.ParseToken(token)
		if err != nil || email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Load admin credentials from .env
		if err := godotenv.Load(".env"); err != nil {
			log.Fatal("Error loading .env file")
		}
		if email != os.Getenv("Admin_Email") {
			c.Redirect(http.StatusSeeOther, "/admin/login")
			c.Abort()
			return
		}

		c.Next()
	}
}
