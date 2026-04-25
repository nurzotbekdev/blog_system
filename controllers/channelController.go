package controllers

import (
	"blog_system/logging"
	"blog_system/schemas"
	"blog_system/services"
	"blog_system/validators"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ChannelController struct {
	ChannelService services.ChannelService
}

func NewChannelController(channel services.ChannelService) *ChannelController {
	return &ChannelController{ChannelService: channel}
}

func (channel *ChannelController) Create(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized create channel attempt")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	req := schemas.CreateChannelRequest{
		UserID:      currentUserID,
		Name:        ctx.PostForm("name"),
		Description: ctx.PostForm("description"),
	}

	profileImage, _ := ctx.FormFile("profile_image")
	bannerImage, _ := ctx.FormFile("banner_image")

	req.ProfileFile = profileImage
	req.BannerFile = bannerImage

	err := validators.ValidateChannel(req.Name, profileImage, bannerImage)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = channel.ChannelService.CreateChannel(req)
	if err != nil {
		if errors.Is(err, services.ErrChannelAlreadyExists) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		logging.Error("create channel failed", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Channel created",
	})
}

func (channel *ChannelController) MyChannel(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized create channel attempt")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	channelData, err := channel.ChannelService.GetMyChannel(currentUserID)
	if err != nil {
		if errors.Is(err, services.ErrChannelNotFound) {
			logging.Log.Warn("Channel not found", zap.Uint("user_id", currentUserID))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Get my channel failed", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("My channel fetched", zap.Uint("user_id", currentUserID))
	ctx.JSON(http.StatusOK, channelData)
}

func (channel *ChannelController) Channel(ctx *gin.Context) {
	name := ctx.Query("name")

	channelData, err := channel.ChannelService.GetChannel(name)
	if err != nil {
		logging.Log.Error("Get my channel failed", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": channelData,
	})
}

func (channel *ChannelController) UpdateChannel(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized create channel attempt")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	nameStr := ctx.PostForm("name")
	descriptionStr := ctx.PostForm("description")

	var name *string
	var description *string

	if nameStr != "" {
		name = &nameStr
	}

	if descriptionStr != "" {
		description = &descriptionStr
	}

	profileImage, _ := ctx.FormFile("profile_image")
	bannerImage, _ := ctx.FormFile("banner_image")

	err := channel.ChannelService.EditChannel(currentUserID, name, description, profileImage, bannerImage)
	if err != nil {
		if errors.Is(err, services.ErrChannelNotFound) {
			logging.Log.Warn("Channel not found for update", zap.Uint("user_id", currentUserID))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, services.ErrChannelAlreadyExists) {
			logging.Log.Warn("Duplicate channel name", zap.Uint("user_id", currentUserID))
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, services.ErrNoFieldsToUpdate) {
			logging.Log.Warn("Empty update request", zap.Uint("user_id", currentUserID))
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Update channel failed", zap.Uint("user_id", currentUserID), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Channel updated successfully", zap.Uint("user_id", currentUserID))
	ctx.JSON(http.StatusOK, gin.H{"message": "Channel updated successfully"})
}
