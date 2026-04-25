package routes

import (
	"blog_system/controllers"
	"blog_system/middleware"
	"blog_system/services"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	userService := services.NewUserServices()
	userController := controllers.NewUserController(userService)

	r.GET("/auth/google/login", userController.GoogleLogin)
	r.GET("/auth/google/callback", userController.GoogleCallback)
	r.GET("/profile", middleware.AuthMiddleware(), userController.MyProfile)
}
