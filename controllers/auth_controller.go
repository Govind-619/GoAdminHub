package controllers

import (
	"github.com/Govind-619/GoAdminHub/utils"

	"github.com/Govind-619/GoAdminHub/models"

	"net/http"
	"regexp"
	"strings"

	"github.com/Govind-619/GoAdminHub/config"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ShowLoginPage renders the login view.
func ShowLoginPage(c *gin.Context) {
	if token, err := c.Cookie("jwtToken"); err == nil && token != "" {
		c.Redirect(http.StatusSeeOther, "/home")
		return
	}
	c.HTML(http.StatusOK, "login.html", nil)
}

// Login processes the login form submission.
func Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.HTML(http.StatusNotFound, "login.html", gin.H{"error": "User not found"})
		return
	}

	// Validate password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Set JWT cookie
	c.SetCookie("jwtToken", token, 36000, "/", "", false, true)
	c.Redirect(http.StatusSeeOther, "/home")
}

// ShowSignupPage renders the signup view.
func ShowSignupPage(c *gin.Context) {
	if token, err := c.Cookie("jwtToken"); err == nil && token != "" {
		c.Redirect(http.StatusSeeOther, "/home")
		return
	}
	c.HTML(http.StatusOK, "signup.html", nil)
}

// Signup processes the signup form submission.
func Signup(c *gin.Context) {
	username := strings.TrimSpace(c.PostForm("name"))
	email := strings.TrimSpace(c.PostForm("email"))
	password := c.PostForm("password")
	cPassword := c.PostForm("c_password")

	// Validate username
	usernameRegex := regexp.MustCompile(`^[a-zA-Z\s]+$`)
	if !usernameRegex.MatchString(username) {
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{"error": "Username should only contain letters"})
		return
	}
	if len(username) < 4 || len(username) > 20 {
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{"error": "Username must be between 4 and 20 characters"})
		return
	}

	// Validate email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{"error": "Invalid email address"})
		return
	}
	if username == "" || email == "" || password == "" {
		c.HTML(http.StatusConflict, "signup.html", gin.H{"error": "Please fill all the details"})
		return
	}

	// Check for existing user
	var existingUser models.User
	if err := config.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		c.HTML(http.StatusConflict, "signup.html", gin.H{"error": "Email already in use"})
		return
	}

	if password != cPassword {
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{"error": "Password mismatch"})
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
	c.HTML(http.StatusOK, "login.html", gin.H{"message": "Signup successful, please login."})
}

// Home renders the user home view.
func Home(c *gin.Context) {
	if token, err := c.Cookie("jwtToken"); err != nil || token == "" {
		c.Redirect(http.StatusSeeOther, "/auth/login")
		return
	}
	email, exists := c.Get("email")
	if !exists {
		c.Redirect(http.StatusSeeOther, "/auth/login")
		return
	}
	c.HTML(http.StatusOK, "home.html", gin.H{"email": email})
}

// Logout clears the JWT cookie and logs the user out.
func Logout(c *gin.Context) {
	c.SetCookie("jwtToken", "", -1, "/", "", false, true)
	c.Redirect(http.StatusSeeOther, "/auth/login")
}
