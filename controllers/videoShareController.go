package controllers

import (
	"blog_system/helper"
	"blog_system/models"
	"blog_system/schemas"
	"blog_system/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
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
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	videoID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var body schemas.VideoShareSchemas
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	newVideoShare := models.VideoShare{
		UserID:   currentUserID,
		VideoID:  videoID,
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
