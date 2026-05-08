package services

import (
	"blog_system/config"
	"blog_system/helper"
	"blog_system/logging"
	"blog_system/models"
	"blog_system/schemas"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type VideoService interface {
	CreateVideo(req schemas.CreateVideoRequest, userID uint) error
}

type videoService struct{}

func NewVideoServices() VideoService {
	return &videoService{}
}

func (s *videoService) CreateVideo(req schemas.CreateVideoRequest, userID uint) error {
	var channel models.Channel
	if err := config.DB.Where("user_id = ?", userID).First(&channel).Error; err != nil {
		return ErrChannelNotFound
	}

	var language models.Language
	if err := config.DB.First(&language, req.LanguageID).Error; err != nil {
		return ErrLanguageNotFound
	}

	var category models.Category
	if err := config.DB.First(&category, req.CategoryID).Error; err != nil {
		return ErrCategoryNotFound
	}

	var videoPath, thumbnailPath string
	var err error
	var size int64

	if req.FilePath != nil {
		videoPath, err = helper.SaveFile(req.FilePath, "uploads/video")
		if err != nil {
			logging.Log.Error("Profile upload failed")
			return nil
		}
	}

	if req.ThumbnailPath != nil {
		thumbnailPath, err = helper.SaveFile(req.ThumbnailPath, "uploads/thumbnail")
		if err != nil {
			logging.Log.Error("Profile upload failed")
			return nil
		}
	}

	if req.FilePath != nil {
		size = req.FilePath.Size
	}

	duration, err := helper.GetVideoDuration(videoPath)
	if err != nil {
		return err
	}

	resolution, err := helper.GetResolution(videoPath)
	if err != nil {
		return err
	}

	tags := helper.GenerateTags(req.Title, req.Description)

	newVideo := models.Video{
		ChannelID:     channel.ID,
		LanguageID:    language.ID,
		CategoryID:    category.ID,
		Title:         req.Title,
		Description:   req.Description,
		FilePath:      videoPath,
		ThumbnailPath: thumbnailPath,
		Tags:          strings.Join(tags, ","),
		Resolution:    resolution,
		Size:          size,
		DurationVideo: duration,
		Visibility:    "public",
	}

	if err := config.DB.Create(&newVideo).Error; err != nil {
		return err
	}

	logging.Log.Info("Video created successfully", zap.Uint("video_id", newVideo.ID))
	if err := config.DB.
		Model(&models.Channel{}).
		Where("id = ?", channel.ID).
		Update("total_videos", gorm.Expr("total_videos + ?", 1)).Error; err != nil {

		logging.Log.Error(
			"Failed to update video count",
			zap.Uint("video_id", newVideo.ID),
			zap.Error(err),
		)

		return err
	}

	logging.Log.Error("Video upload failed", zap.Error(err))
	return err
}
