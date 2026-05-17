package services

import (
	"blog_system/config"
	"blog_system/helper"
	"blog_system/jobs"
	"blog_system/logging"
	"blog_system/models"
	"blog_system/schemas"
	"encoding/json"
	"errors"
	"mime/multipart"
	"os"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type VideoService interface {
	CreateVideo(req schemas.CreateVideoRequest, userID uint) error
	GetMyVideo(userID uint) ([]schemas.MyVideoResponse, error)
	GetVideo(page, limit int, categoryID uint, search, languageCode, sortBy string) (*schemas.VideoListResponse, error)
	GetVideoByID(videoID, userID uint) (*schemas.VideoResponse, error)
	EditVideo(videoID, userID uint, languageID, categoryID *uint, title, description, visibility *string, thumbnailPath *multipart.FileHeader) error
	DeleteVideo(videoID, userID uint) error
}

type videoService struct{}

func NewVideoServices() VideoService {
	return &videoService{}
}

var (
	ErrVideoNotFound = errors.New("Video not found")
)

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
			return err
		}
	}

	if req.ThumbnailPath != nil {
		thumbnailPath, err = helper.SaveFile(req.ThumbnailPath, "uploads/thumbnail")
		if err != nil {
			logging.Log.Error("Profile upload failed")
			return err
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
		Status:        "processing",
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

	job := jobs.VideoJob{
		VideoID:    newVideo.ID,
		FilePath:   videoPath,
		Resolution: resolution,
	}

	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	config.RedisClient.LPush(config.Ctx, "video_queue", data)

	return nil
}

func (s *videoService) GetMyVideo(userID uint) ([]schemas.MyVideoResponse, error) {
	var channel models.Channel
	if err := config.DB.Where("user_id = ?", userID).First(&channel).Error; err != nil {
		return nil, ErrChannelNotFound
	}

	var results []schemas.MyVideoResponse
	tx := config.DB.Table("videos").Select(`
		videos.id as video_id,
		channels.name as channel_name,
		channels.profile_image,
		languages.name as language_name,
		languages.code,
		categories.name as category_name,
		videos.title,
		videos.description,
		videos.file_path,
		videos.thumbnail_path,
		videos.tags,
		videos.resolution,
		videos.size,
		videos.views,
		videos.like_count,
		videos.comment_count,
		videos.dislike_count,
		videos.duration_video,
		videos.visibility,
		videos.status,
		videos.share_count,
		videos.download_count,
		videos.created_at
		`).
		Joins("JOIN channels ON channels.id = videos.channel_id").
		Joins("JOIN languages ON languages.id = videos.language_id").
		Joins("JOIN categories ON categories.id = videos.category_id").
		Where("videos.channel_id = ?", channel.ID).
		Scan(&results)

	if tx.Error != nil {
		return nil, tx.Error
	}

	if tx.RowsAffected == 0 {
		return nil, ErrVideoNotFound
	}

	for i := range results {
		var qualities []models.VideoQuality
		if err := config.DB.
			Where("video_id = ?", results[i].VideoID).
			Find(&qualities).Error; err != nil {
			continue
		}

		for _, q := range qualities {
			results[i].Qualities = append(
				results[i].Qualities,
				schemas.VideoQualityResponse{
					Quality:  q.Quality,
					VideoURL: q.VideoURL,
					Size:     q.Size,
					Format:   q.Format,
				},
			)
		}
	}

	return results, nil
}

func (s *videoService) GetVideo(page, limit int, categoryID uint, search, languageCode, sortBy string) (*schemas.VideoListResponse, error) {
	var results []schemas.VideoResponse
	var total int64

	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	if limit > 50 {
		limit = 50
	}

	offset := (page - 1) * limit

	query := config.DB.Table("videos").
		Select(`
			videos.id as video_id,
			channels.name as channel_name,
			channels.profile_image,
			languages.name as language_name,
			languages.code,
			categories.name as category_name,
			videos.title,
			videos.description,
			videos.file_path,
			videos.thumbnail_path,
			videos.tags,
			videos.resolution,
			videos.size,
			videos.views,
			videos.like_count,
			videos.duration_video,
			videos.created_at
		`).
		Joins("JOIN channels ON channels.id = videos.channel_id").
		Joins("JOIN languages ON languages.id = videos.language_id").
		Joins("JOIN categories ON categories.id = videos.category_id").
		Where("videos.visibility = ?", "public")

	if search != "" {
		query = query.Where(
			"videos.title ILIKE ? OR videos.tags ILIKE ?",
			"%"+search+"%",
			"%"+search+"%",
		)
	}

	if categoryID != 0 {
		query = query.Where("videos.category_id = ?", categoryID)
	}

	if languageCode != "" {
		query = query.Where("languages.code = ?", languageCode)
	}

	countQuery := query

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	switch sortBy {
	case "trending":
		query = query.Order("videos.views DESC")

	case "popular":
		query = query.Order("videos.like_count DESC")

	case "latest":
		fallthrough

	default:
		query = query.Order("videos.created_at DESC")
	}

	tx := query.
		Limit(limit).
		Offset(offset).
		Scan(&results)

	if tx.Error != nil {
		return nil, tx.Error
	}

	for i := range results {
		var qualities []models.VideoQuality

		if err := config.DB.
			Where("video_id = ?", results[i].VideoID).
			Find(&qualities).Error; err != nil {
			continue
		}

		for _, q := range qualities {
			results[i].Qualities = append(
				results[i].Qualities,
				schemas.VideoQualityResponse{
					Quality:  q.Quality,
					VideoURL: q.VideoURL,
					Size:     q.Size,
					Format:   q.Format,
				},
			)
		}
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	response := schemas.VideoListResponse{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		Data:       results,
	}

	return &response, nil
}

func (s *videoService) GetVideoByID(videoID, userID uint) (*schemas.VideoResponse, error) {
	var results schemas.VideoResponse

	tx := config.DB.Table("videos").
		Select(`
			videos.id as video_id,
			channels.name as channel_name,
			channels.profile_image,
			languages.name as language_name,
			languages.code,
			categories.name as category_name,
			videos.title,
			videos.description,
			videos.file_path,
			videos.thumbnail_path,
			videos.tags,
			videos.resolution,
			videos.size,
			videos.views,
			videos.like_count,
			videos.comment_count,
			videos.dislike_count,
			videos.duration_video,
			videos.visibility,
			videos.share_count,
			videos.download_count,
			videos.favorite_count,
			videos.created_at
	`).
		Joins("JOIN channels ON channels.id = videos.channel_id").
		Joins("JOIN languages ON languages.id = videos.language_id").
		Joins("JOIN categories ON categories.id = videos.category_id").
		Where("videos.id = ? AND videos.visibility = ?", videoID, "public").
		Scan(&results)

	if tx.Error != nil {
		return nil, tx.Error
	}

	if tx.RowsAffected == 0 {
		return nil, ErrVideoNotFound
	}

	var qualities []models.VideoQuality
	if err := config.DB.
		Where("video_id = ?", videoID).
		Find(&qualities).Error; err != nil {
		return nil, err
	}

	for _, q := range qualities {

		results.Qualities = append(
			results.Qualities,
			schemas.VideoQualityResponse{
				Quality:  q.Quality,
				VideoURL: q.VideoURL,
				Size:     q.Size,
				Format:   q.Format,
			},
		)
	}

	if userID != 0 {

		job := jobs.VideoViewJob{
			UserID:  userID,
			VideoID: videoID,
		}

		data, err := json.Marshal(job)
		if err == nil {
			_ = config.RedisClient.LPush(
				config.Ctx,
				"video_view_queue",
				data,
			).Err()
		}
	}

	return &results, nil
}

func (s *videoService) EditVideo(videoID, userID uint, languageID, categoryID *uint, title, description, visibility *string, thumbnailPath *multipart.FileHeader) error {
	var channel models.Channel
	if err := config.DB.Where("user_id = ?", userID).First(&channel).Error; err != nil {
		return ErrChannelNotFound
	}

	var video models.Video
	if err := config.DB.Where("id = ? AND channel_id = ?", videoID, channel.ID).First(&video).Error; err != nil {
		return ErrVideoNotFound
	}

	updates := map[string]interface{}{}
	var oldThumbnailPath string

	if languageID != nil {
		var language models.Language
		if err := config.DB.First(&language, languageID).Error; err != nil {
			return ErrLanguageNotFound
		}

		updates["language_id"] = *languageID
	}

	if categoryID != nil {
		var category models.Category
		if err := config.DB.First(&category, categoryID).Error; err != nil {
			return ErrCategoryNotFound
		}

		updates["category_id"] = *categoryID
	}

	if title != nil {
		updates["title"] = *title
	}

	if description != nil {
		updates["description"] = *description
	}

	if visibility != nil {
		if *visibility != "public" && *visibility != "private" {
			return errors.New("invalid visibility")
		}
		updates["visibility"] = *visibility
	}

	if thumbnailPath != nil {
		oldThumbnailPath = video.ThumbnailPath

		ThumbnailPath, err := helper.SaveFile(thumbnailPath, "uploads/thumbnail")
		if err != nil {
			return err
		}

		updates["thumbnail_path"] = ThumbnailPath
	}

	if title != nil || description != nil {

		newTitle := video.Title
		newDesc := video.Description

		if title != nil {
			newTitle = *title
		}

		if description != nil {
			newDesc = *description
		}

		tags := helper.GenerateTags(newTitle, newDesc)
		updates["tags"] = strings.Join(tags, ",")
	}

	if len(updates) == 0 {
		return ErrNoFieldsToUpdate
	}

	if err := config.DB.Model(&video).Updates(updates).Error; err != nil {
		return err
	}

	if oldThumbnailPath != "" {
		if err := os.Remove(oldThumbnailPath); err != nil {
			logging.Log.Warn("Failed to delete old thumbnail image", zap.String("path", oldThumbnailPath), zap.Error(err))
		} else {
			logging.Log.Info("Old profile image deleted")
		}
	}

	return nil
}

func (s *videoService) DeleteVideo(videoID, userID uint) error {
	var channel models.Channel
	if err := config.DB.Where("user_id = ?", userID).First(&channel).Error; err != nil {
		return ErrChannelNotFound
	}

	var video models.Video
	if err := config.DB.Where("id = ? AND channel_id = ?", videoID, channel.ID).First(&video).Error; err != nil {
		return ErrVideoNotFound
	}

	if err := config.DB.Delete(&video).Error; err != nil {
		return err
	}

	if video.ThumbnailPath != "" {
		_ = os.Remove(video.ThumbnailPath)
	}

	if video.FilePath != "" {
		_ = os.Remove(video.FilePath)
	}

	if err := config.DB.
		Model(&models.Channel{}).
		Where("id = ?", channel.ID).
		UpdateColumn("total_videos",
			gorm.Expr("GREATEST(total_videos - 1,0)")).Error; err != nil {
		return err
	}

	return nil
}
