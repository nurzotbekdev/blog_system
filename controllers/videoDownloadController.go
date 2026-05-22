package controllers

import (
	"blog_system/helper"
	"blog_system/logging"
	"blog_system/models"
	"blog_system/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
		logging.Log.Warn("Unauthorized video download attempt", zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	videoID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		logging.Log.Warn("Invalid video id param for download", zap.Uint("user_id", currentUserID), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newDownload := models.VideoDownload{
		UserID:  currentUserID,
		VideoID: videoID,
	}

	logging.Log.Info("Video download request started", zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID))
	if err := download.VideoDownloadService.CreateVideoDownload(newDownload); err != nil {
		if errors.Is(err, services.ErrVideoNotFound) {
			logging.Log.Warn("Video not found for download", zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID), zap.Error(err))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Failed to create video download", zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Video downloaded successfully", zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID))
	ctx.JSON(http.StatusCreated, gin.H{"message": "Video download successful"})
}
