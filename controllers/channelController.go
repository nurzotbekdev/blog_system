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

	logging.Log.Info("Creating channel", zap.Uint("user_id", currentUserID), zap.String("channel_name", req.Name))
	err := validators.ValidateChannel(req.Name, profileImage, bannerImage)
	if err != nil {
		logging.Log.Warn("Channel validation failed", zap.Uint("user_id", currentUserID), zap.String("channel_name", req.Name), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = channel.ChannelService.CreateChannel(req)
	if err != nil {
		logging.Log.Warn("Channel already exists", zap.Uint("user_id", currentUserID), zap.String("channel_name", req.Name))
		if errors.Is(err, services.ErrChannelAlreadyExists) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Create channel failed", zap.Uint("user_id", currentUserID), zap.String("channel_name", req.Name), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Channel created successfully", zap.Uint("user_id", currentUserID))
	ctx.JSON(http.StatusCreated, gin.H{"message": "Channel created"})
}

func (channel *ChannelController) MyChannel(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized my channel access")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	logging.Log.Info("Fetching my channel", zap.Uint("user_id", currentUserID))
	channelData, err := channel.ChannelService.GetMyChannel(currentUserID)
	if err != nil {
		if errors.Is(err, services.ErrChannelNotFound) {
			logging.Log.Warn("Channel not found", zap.Uint("user_id", currentUserID))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Get my channel failed", zap.Uint("user_id", currentUserID), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("My channel fetched successfully", zap.Uint("user_id", currentUserID), zap.Uint("channel_id", channelData.ID), zap.String("channel_name", channelData.Name))
	ctx.JSON(http.StatusOK, channelData)
}

func (channel *ChannelController) Channel(ctx *gin.Context) {
	name := ctx.Query("name")
	if name == "" {
		logging.Log.Warn("Empty channel name query")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Channel name is required"})
		return
	}

	logging.Log.Info("Fetching public channel", zap.String("channel_name", name))
	channelData, err := channel.ChannelService.GetChannel(name)
	if err != nil {
		logging.Log.Error("Get public channel failed", zap.String("channel_name", name), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Public channel fetched successfully")
	ctx.JSON(http.StatusOK, gin.H{"data": channelData})
}

func (channel *ChannelController) UpdateChannel(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized update channel attempt")
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

	logging.Log.Info("Updating channel", zap.Uint("user_id", currentUserID), zap.String("new_name", nameStr))
	if err := channel.ChannelService.EditChannel(currentUserID, name, description, profileImage, bannerImage); err != nil {
		if errors.Is(err, services.ErrChannelNotFound) {
			logging.Log.Warn("Channel not found for update", zap.Uint("user_id", currentUserID))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, services.ErrChannelAlreadyExists) {
			logging.Log.Warn("Duplicate channel name", zap.Uint("user_id", currentUserID), zap.String("channel_name", nameStr))
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, services.ErrNoFieldsToUpdate) {
			logging.Log.Warn("Empty update request", zap.Uint("user_id", currentUserID))
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Update channel failed", zap.Uint("user_id", currentUserID), zap.String("channel_name", nameStr), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Channel updated successfully", zap.Uint("user_id", currentUserID), zap.String("channel_name", nameStr))
	ctx.JSON(http.StatusOK, gin.H{"message": "Channel updated successfully"})
}
