package routes

import (
	"blog_system/controllers"
	"blog_system/middleware"
	"blog_system/services"

	"github.com/gin-gonic/gin"
)

func ChannelRoutes(r *gin.Engine) {
	channelService := services.NewChannelService()
	channelController := controllers.NewChannelController(channelService)

	r.POST("/channel", middleware.AuthMiddleware(), channelController.Create)
	r.GET("/channel/me", middleware.AuthMiddleware(), channelController.MyChannel)
	r.GET("/channel", middleware.AuthMiddleware(), channelController.Channel)
	r.PUT("/channel", middleware.AuthMiddleware(), channelController.UpdateChannel)
}
