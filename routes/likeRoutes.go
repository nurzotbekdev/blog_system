package routes

import (
	"blog_system/controllers"
	"blog_system/middleware"
	"blog_system/services"

	"github.com/gin-gonic/gin"
)

func LikeRoutes(r *gin.Engine) {
	likeService := services.NewLikeServices()
	likeController := controllers.NewLikeController(likeService)

	r.POST("/like", middleware.AuthMiddleware(), likeController.Create)
	r.PUT("/video/:id/like", middleware.AuthMiddleware(), likeController.UpdateLike)
	r.DELETE("/video/:id/like", middleware.AuthMiddleware(), likeController.DeleteLike)
	r.GET("/video/:id/reaction", middleware.AuthMiddleware(), likeController.GetReaction)
	r.GET("/user/likes", middleware.AuthMiddleware(), likeController.GetUserLikes)
}
