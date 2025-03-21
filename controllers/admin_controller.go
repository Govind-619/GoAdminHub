package controllers

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/Govind-619/GoAdminHub/utils"

	"github.com/Govind-619/GoAdminHub/models"

	"github.com/Govind-619/GoAdminHub/config"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

// ShowAdminLogin renders the admin login view.
func ShowAdminLogin(c *gin.Context) {
	if token, err := c.Cookie("Admin"); err == nil && token != "" {
		c.Redirect(http.StatusSeeOther, "/admin/panel")
		return
	}
	c.HTML(http.StatusOK, "admin-login.html", nil)
}

// AdminLogin processes the admin login form.
func AdminLogin(c *gin.Context) {
	username := c.PostForm("admin")
	password := c.PostForm("password")

	// Load admin credentials from .env
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	if username != os.Getenv("Admin_Email") || password != os.Getenv("Admin_Password") {
		c.HTML(http.StatusUnauthorized, "admin-login.html", gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateToken(username, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.SetCookie("Admin", token, 36000, "/", "", false, true)
	c.Redirect(http.StatusSeeOther, "/admin/panel")
}

// AdminPanel renders the admin panel view with a list of users.
func AdminPanel(c *gin.Context) {
	var users []models.User
	if err := config.DB.Find(&users).Error; err != nil {
		c.Redirect(http.StatusSeeOther, "/admin/login")
		return
	}
	successMessage := c.Query("success")
	errorMessage := c.Query("error")
	c.HTML(http.StatusOK, "adminpanel.html", gin.H{
		"users":   users,
		"success": successMessage,
		"error":   errorMessage,
	})
}

// AdminLogout logs the admin out.
func AdminLogout(c *gin.Context) {
	c.SetCookie("Admin", "", -1, "/", "", false, true)
	c.Redirect(http.StatusSeeOther, "/admin/login")
}

// AddUser allows the admin to add a new user.
func AddUser(c *gin.Context) {
	username := strings.TrimSpace(c.PostForm("name"))
	email := strings.TrimSpace(c.PostForm("email"))
	password := c.PostForm("password")

	if !isValidUsername(username) {
		c.Redirect(http.StatusSeeOther, "/admin/panel?error=Invalid+username")
		return
	}
	if !isValidEmail(email) {
		c.Redirect(http.StatusSeeOther, "/admin/panel?error=Invalid+email")
		return
	}

	var existingUser models.User
	if err := config.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		c.Redirect(http.StatusSeeOther, "/admin/panel?error=User+already+exists")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	newUser := models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}
	config.DB.Create(&newUser)
	c.Redirect(http.StatusSeeOther, "/admin/panel?success=User+added+successfully")
}

// SearchUser allows the admin to search for users by username.
func SearchUser(c *gin.Context) {
	query := c.Query("query")
	var users []models.User
	if err := config.DB.Where("username ILIKE ?", "%"+query+"%").Find(&users).Error; err != nil {
		c.Redirect(http.StatusSeeOther, "/admin/panel")
		return
	}
	c.HTML(http.StatusOK, "adminpanel.html", gin.H{"users": users})
}

// EditUser allows the admin to edit a user’s details.
func EditUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Username = c.PostForm("name")
	user.Email = c.PostForm("email")
	newPassword := c.PostForm("password")
	if newPassword != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		user.Password = string(hashedPassword)
	}

	config.DB.Save(&user)
	c.Redirect(http.StatusSeeOther, "/admin/panel")
}

// DeleteUser allows the admin to delete a user.
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "login.html", gin.H{"error": "User not found"})
		return
	}
	// If the admin is deleting their own account, log them out.
	if email, exists := c.Get("email"); exists && user.Email == email {
		c.SetCookie("jwtToken", "", -1, "/", "", false, true)
		c.Redirect(http.StatusSeeOther, "/auth/login")
		return
	}
	config.DB.Delete(&user)
	c.Redirect(http.StatusSeeOther, "/admin/panel")
}

// Helper functions to validate username and email.
func isValidUsername(username string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z\s]{4,20}$`)
	return regex.MatchString(username)
}

func isValidEmail(email string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return regex.MatchString(email)
}
