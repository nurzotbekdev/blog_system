package routes

import (
	"blog_system/controllers"
	"blog_system/middleware"
	"blog_system/services"

	"github.com/gin-gonic/gin"
)

func PlaylistRoutes(r *gin.Engine) {
	playlistService := services.NewPlaylistServices()
	playlistController := controllers.NewPlaylistController(playlistService)

	r.POST("/playlist", middleware.AuthMiddleware(), playlistController.Create)
	r.GET("/playlist/me", middleware.AuthMiddleware(), playlistController.GetMyPlaylist)
	r.GET("/playlist", middleware.AuthMiddleware(), playlistController.GetPlaylist)
	r.PUT("/plalist/:id", middleware.AuthMiddleware(), playlistController.UpdatePlaylist)
	r.DELETE("/playlist/:id", middleware.AuthMiddleware(), playlistController.DeletePlaylist)
}
