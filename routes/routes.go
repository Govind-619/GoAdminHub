package routes

import (
	"net/http"

	"github.com/Govind-619/GoAdminHub/middlewares"

	"github.com/Govind-619/GoAdminHub/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Serve static files and load templates
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*.html")
	r.Use(middlewares.ClearCache())

	// In routes/routes.go, before setting up other groups:
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/auth/login")
	})

	// Public routes for user authentication
	auth := r.Group("/auth")
	{
		auth.GET("/login", controllers.ShowLoginPage)
		auth.POST("/login", controllers.Login)
		auth.GET("/signup", controllers.ShowSignupPage)
		auth.POST("/signup", controllers.Signup)
	}

	// Protected routes for authenticated users
	userRoutes := r.Group("/")
	userRoutes.Use(middlewares.Authenticate())
	{
		userRoutes.GET("/home", controllers.Home)
		userRoutes.GET("/logout", controllers.Logout)
	}

	// Admin public login routes with no-cache middleware applied
	r.GET("/admin/login", middlewares.ClearCache(), controllers.ShowAdminLogin)
	r.POST("/admin/login", middlewares.ClearCache(), controllers.AdminLogin)

	// Protected routes for admin actions
	adminRoutes := r.Group("/admin")
	adminRoutes.Use(middlewares.AdminAuthenticate())
	{
		adminRoutes.GET("/panel", controllers.AdminPanel)
		adminRoutes.GET("/logout", controllers.AdminLogout)
		adminRoutes.POST("/adduser", controllers.AddUser)
		adminRoutes.GET("/search", controllers.SearchUser)
		adminRoutes.POST("/edituser/:id", controllers.EditUser)
		adminRoutes.GET("/deleteuser/:id", controllers.DeleteUser)
	}
}
