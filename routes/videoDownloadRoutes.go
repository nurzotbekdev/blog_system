package routes

import (
	"blog_system/controllers"
	"blog_system/middleware"
	"blog_system/services"

	"github.com/gin-gonic/gin"
)

func VideoDownloadRoutes(r *gin.Engine) {
	videoDownloadService := services.NewVideoDownloadServices()
	videoDownloadController := controllers.NewVideoDownloadController(videoDownloadService)

	r.POST("/video/:id/download", middleware.AuthMiddleware(), videoDownloadController.Create)
}
