package controllers

import (
	"blog_system/helper"
	"blog_system/logging"
	"blog_system/models"
	"blog_system/schemas"
	"blog_system/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type VideoShareController struct {
	VideoShareService services.VideoShareService
}

func NewVideoShareController(share services.VideoShareService) *VideoShareController {
	return &VideoShareController{VideoShareService: share}
}

func (share *VideoShareController) Create(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized create channel attempt")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	VideoID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		logging.Log.Warn("Invalid subscription id param", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var body schemas.VideoShareSchemas
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logging.Log.Warn("Invalid JSON format", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	newVideoShare := models.VideoShare{
		UserID:   currentUserID,
		VideoID:  VideoID,
		Platform: body.Platform,
	}

	if err := share.VideoShareService.CreateVideoShare(newVideoShare); err != nil {
		if errors.Is(err, services.ErrVideoNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, services.ErrVideoNotVisibility) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Video share successful"})
}
