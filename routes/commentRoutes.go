package routes

import (
	"blog_system/controllers"
	"blog_system/middleware"
	"blog_system/services"

	"github.com/gin-gonic/gin"
)

func CommentRoutes(r *gin.Engine) {
	commentService := services.NewCommentServices()
	commentController := controllers.NewCommentController(commentService)

	r.POST("/comment", middleware.AuthMiddleware(), commentController.Create)
	r.GET("/video/:id/comments", middleware.AuthMiddleware(), commentController.GetVideoComments)
}
