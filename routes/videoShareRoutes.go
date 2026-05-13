package routes

import (
	"blog_system/controllers"
	"blog_system/middleware"
	"blog_system/services"

	"github.com/gin-gonic/gin"
)

func VideoShareRoutes(r *gin.Engine) {
	videoShareService := services.NewVideoShareServices()
	videoShareController := controllers.NewVideoShareController(videoShareService)

	r.POST("/video/:id/share", middleware.AuthMiddleware(), videoShareController.Create)
}
