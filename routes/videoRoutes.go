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
	r.GET("/video/me", middleware.AuthMiddleware(), videoController.MyVideo)
	r.GET("/video", middleware.AuthMiddleware(), videoController.ListVideos)
	r.GET("/video/:id", middleware.AuthMiddleware(), videoController.VideoDetail)
	r.PUT("/video/:id", middleware.AuthMiddleware(), videoController.UpdateVideo)
	r.DELETE("/video/:id", middleware.AuthMiddleware(), videoController.DeleteVideo)
}
