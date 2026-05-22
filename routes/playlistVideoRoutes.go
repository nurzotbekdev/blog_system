package routes

import (
	"blog_system/controllers"
	"blog_system/middleware"
	"blog_system/services"

	"github.com/gin-gonic/gin"
)

func PlaylistVideoRoutes(r *gin.Engine) {
	playlistVideoService := services.NewPlaylistVideoServices()
	playlistVideoController := controllers.NewPlaylistVideoController(playlistVideoService)

	r.POST("/playlist/video", middleware.AuthMiddleware(), playlistVideoController.Create)
	r.GET("/playlist/:id/videos", middleware.AuthMiddleware(), playlistVideoController.GetPlaylistDetail)
	r.GET("/playlist/:id", middleware.AuthMiddleware(), playlistVideoController.GetMyPlaylistDetail)
	r.DELETE("/playlist/video/:id", middleware.AuthMiddleware(), playlistVideoController.DeletePlaylistVideo)
}
