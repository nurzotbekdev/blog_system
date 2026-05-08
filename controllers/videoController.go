package controllers

import (
	"blog_system/logging"
	"blog_system/schemas"
	"blog_system/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type VideoController struct {
	VideoService services.VideoService
}

func NewVideoController(video services.VideoService) *VideoController {
	return &VideoController{VideoService: video}
}

func (video *VideoController) Create(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized create video attempt")

		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	currentUserID := userID.(uint)

	var req schemas.CreateVideoRequest

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format"})
		return
	}

	filePath, _ := ctx.FormFile("file_path")
	thumbnailPath, _ := ctx.FormFile("thumbnail_path")

	req.FilePath = filePath
	req.ThumbnailPath = thumbnailPath

	if err := video.VideoService.CreateVideo(req, currentUserID); err != nil {

		if errors.Is(err, services.ErrChannelNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		if errors.Is(err, services.ErrLanguageNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		if errors.Is(err, services.ErrCategoryNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		logging.Log.Error("create video failed")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Video created successfully",
	})
}
