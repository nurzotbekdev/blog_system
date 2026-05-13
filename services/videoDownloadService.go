package services

import (
	"blog_system/config"
	"blog_system/models"

	"gorm.io/gorm"
)

type VideoDownloadService interface {
	CreateVideoDownload(download models.VideoDownload) error
}

type videoDownloadService struct{}

func NewVideoDownloadServices() VideoDownloadService {
	return &videoDownloadService{}
}

func (s *videoDownloadService) CreateVideoDownload(download models.VideoDownload) error {
	var video models.Video

	if err := config.DB.
		Where("id = ?", download.VideoID).
		First(&video).Error; err != nil {
		return ErrVideoNotFound
	}

	if video.Visibility != "public" {
		return ErrVideoNotVisibility
	}

	if err := config.DB.Create(&download).Error; err != nil {
		return err
	}

	if err := config.DB.
		Model(&models.Video{}).
		Where("id = ?", video.ID).
		Update("download_count", gorm.Expr("download_count + ?", 1)).
		Error; err != nil {
		return err
	}

	return nil
}
