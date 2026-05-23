package controllers

import (
	"blog_system/helper"
	"blog_system/models"
	"blog_system/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type VideoDownloadController struct {
	VideoDownloadService services.VideoDownloadService
}

func NewVideoDownloadController(download services.VideoDownloadService) *VideoDownloadController {
	return &VideoDownloadController{VideoDownloadService: download}
}

func (download *VideoDownloadController) Create(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	videoID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newDownload := models.VideoDownload{
		UserID:  currentUserID,
		VideoID: videoID,
	}

	if err := download.VideoDownloadService.CreateVideoDownload(newDownload); err != nil {
		if errors.Is(err, services.ErrVideoNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Video download successful"})
}
