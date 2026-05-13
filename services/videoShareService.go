package services

import (
	"blog_system/config"
	"blog_system/models"
	"errors"

	"gorm.io/gorm"
)

type VideoShareService interface {
	CreateVideoShare(share models.VideoShare) error
}

type videoShareService struct{}

func NewVideoShareServices() VideoShareService {
	return &videoShareService{}
}

var (
	ErrVideoNotVisibility = errors.New("Video is not shareable")
)

func (s *videoShareService) CreateVideoShare(share models.VideoShare) error {
	var video models.Video

	if err := config.DB.
		Where("id = ?", share.VideoID).
		First(&video).Error; err != nil {
		return ErrVideoNotFound
	}

	if video.Visibility != "public" {
		return ErrVideoNotVisibility
	}

	if err := config.DB.Create(&share).Error; err != nil {
		return err
	}

	if err := config.DB.
		Model(&models.Video{}).
		Where("id = ?", video.ID).
		Update("share_count", gorm.Expr("share_count + ?", 1)).
		Error; err != nil {
		return err
	}

	return nil
}
