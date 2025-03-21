package main

import (
	"github.com/Govind-619/GoAdminHub/routes"

	"github.com/Govind-619/GoAdminHub/config"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Gin router
	r := gin.Default()

	// Initialize the database connection
	config.InitDB()

	// Set up session store (cookie-based)
	store := cookie.NewStore([]byte("1011"))
	r.Use(sessions.Sessions("login-session", store))

	// Set up routes
	routes.SetupRoutes(r)

	// Run the server on localhost:8080
	r.Run("localhost:8080")
}
