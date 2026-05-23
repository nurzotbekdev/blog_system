package services

import (
	"blog_system/config"
	"blog_system/helper"
	"blog_system/models"
	"blog_system/schemas"
	"errors"
	"mime/multipart"
	"os"
)

type ChannelService interface {
	CreateChannel(req schemas.CreateChannelRequest) (*models.Channel, error)
	GetMyChannel(userID uint) (*schemas.ChannelResponse, error)
	GetChannel(name string) ([]schemas.ChannelSearchResponse, error)
	EditChannel(userID uint, name, description *string, profileImage, bannerImage *multipart.FileHeader) error
}

type channelService struct{}

func NewChannelService() ChannelService {
	return &channelService{}
}

var (
	ErrChannelAlreadyExists = errors.New("Channel already exists")
	ErrChannelNotFound      = errors.New("Channel not found")
	ErrNoFieldsToUpdate     = errors.New("No fields to update")
)

func (s *channelService) CreateChannel(req schemas.CreateChannelRequest) (*models.Channel, error) {
	tx := config.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	var profilePath, bannerPath string
	var err error

	if req.ProfileFile != nil {
		profilePath, err = helper.SaveFile(req.ProfileFile, "uploads/profile")
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if req.BannerFile != nil {
		bannerPath, err = helper.SaveFile(req.BannerFile, "uploads/banner")
		if err != nil {
			helper.RemoveFile(profilePath)
			tx.Rollback()
			return nil, err
		}
	}

	channel := models.Channel{
		UserID:       req.UserID,
		Name:         req.Name,
		Description:  req.Description,
		ProfileImage: profilePath,
		BannerImage:  bannerPath,
	}

	if err := tx.Create(&channel).Error; err != nil {
		helper.RemoveFile(profilePath, bannerPath)
		tx.Rollback()

		if helper.IsDuplicateError(err) {
			return nil, ErrChannelAlreadyExists
		}

		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		helper.RemoveFile(profilePath, bannerPath)
		return nil, err
	}

	return &channel, nil
}

func (s *channelService) GetMyChannel(userID uint) (*schemas.ChannelResponse, error) {
	var result schemas.ChannelResponse
	tx := config.DB.Table("channels").
		Select(`
			channels.id,
			channels.user_id,
			users.email,
			users.full_name,
			users.profile_image as avatar_image,
			channels.name,
			channels.description,
			channels.profile_image,
			channels.banner_image,
			channels.total_subscribers,
			channels.total_videos,
			channels.total_comments,
			channels.total_views,
			channels.created_at
		`).
		Joins("JOIN users ON users.id = channels.user_id").
		Where("channels.user_id = ?", userID).
		Scan(&result)

	if tx.Error != nil {
		return nil, tx.Error
	}

	if tx.RowsAffected == 0 {
		return nil, ErrChannelNotFound
	}

	return &result, nil
}

func (s *channelService) GetChannel(name string) ([]schemas.ChannelSearchResponse, error) {
	var results []schemas.ChannelSearchResponse

	tx := config.DB.Table("channels").
		Select(`
			channels.id,
			channels.name,
			channels.description,
			users.profile_image as avatar_image,
			channels.profile_image,
			channels.banner_image,
			channels.total_subscribers,
			channels.total_videos,
			channels.total_views
		`).
		Joins("JOIN users ON users.id = channels.user_id").
		Where("channels.name ILIKE ?", "%"+name+"%").
		Scan(&results)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return results, nil
}

func (s *channelService) EditChannel(userID uint, name, description *string, profileImage, bannerImage *multipart.FileHeader) error {
	var channel models.Channel
	if err := config.DB.Where("user_id = ?", userID).First(&channel).Error; err != nil {
		return ErrChannelNotFound
	}

	updates := map[string]interface{}{}
	var oldProfileImage, oldBannerImage string

	if name != nil {
		var exists models.Channel
		if err := config.DB.
			Where("name = ? AND user_id != ?", *name, userID).
			First(&exists).Error; err == nil {

			return ErrChannelAlreadyExists
		}

		updates["name"] = *name
	}

	if description != nil {
		updates["description"] = *description
	}

	if profileImage != nil {
		oldProfileImage = channel.ProfileImage

		profilePath, err := helper.SaveFile(profileImage, "uploads/profile")
		if err != nil {
			return err
		}

		updates["profile_image"] = profilePath
	}

	if bannerImage != nil {
		oldBannerImage = channel.BannerImage
		bannerPath, err := helper.SaveFile(bannerImage, "uploads/banner")
		if err != nil {
			return err
		}

		updates["banner_image"] = bannerPath
	}

	if len(updates) == 0 {
		return ErrNoFieldsToUpdate
	}

	if err := config.DB.Model(&channel).Updates(updates).Error; err != nil {
		return err
	}

	if oldProfileImage != "" {
		if err := os.Remove(oldProfileImage); err != nil {
		}
	}

	if oldBannerImage != "" {
		if err := os.Remove(oldBannerImage); err != nil {
		}
	}

	return nil
}
