package routes

import (
	"blog_system/controllers"
	"blog_system/middleware"
	"blog_system/services"

	"github.com/gin-gonic/gin"
)

func VideoRoutes(r *gin.Engine) {
	videoService := services.NewVideoServices()
	videoController := controllers.NewVideoController(videoService)

	r.POST("/video", middleware.AuthMiddleware(), videoController.Create)
}
