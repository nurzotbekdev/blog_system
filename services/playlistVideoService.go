package services

import (
	"blog_system/config"
	"blog_system/models"
	"blog_system/schemas"
	"errors"

	"gorm.io/gorm"
)

type PlaylistVideoService interface {
	CreatePlaylistVideo(playlistVideo models.PlaylistVideo) error
	GetPlaylistDetail(playlistID uint) (*schemas.PlaylistDetailResponse, error)
	GetMyPlaylistVideo(playlistID, userID uint) (*schemas.MyPlaylistDetailResponse, error)
	DeletePlaylistVideo(ID, userID uint) error
}

type playlistVideoService struct{}

func NewPlaylistVideoServices() PlaylistVideoService {
	return &playlistVideoService{}
}

var (
	ErrPlaylistVideoNotFound        = errors.New("Playlist video not found")
	ErrVideoAlreadyExistsInPlaylist = errors.New("video already exists in playlist")
)

func (s *playlistVideoService) CreatePlaylistVideo(playlistVideo models.PlaylistVideo) error {
	var playlist models.Playlist
	if err := config.DB.First(&playlist, playlistVideo.PlaylistID).Error; err != nil {
		return ErrPlaylistNotFound
	}

	var video models.Video
	if err := config.DB.First(&video, playlistVideo.VideoID).Error; err != nil {
		return ErrVideoNotFound
	}

	var exists models.PlaylistVideo

	err := config.DB.
		Where("playlist_id = ? AND video_id = ?",
			playlistVideo.PlaylistID,
			playlistVideo.VideoID,
		).
		First(&exists).Error

	if err == nil {
		return ErrVideoAlreadyExistsInPlaylist
	}

	if err := config.DB.Create(&playlistVideo).Error; err != nil {
		return err
	}

	if err := config.DB.
		Model(&models.Playlist{}).
		Where("id = ?", playlistVideo.PlaylistID).
		UpdateColumn("video_count", gorm.Expr("video_count + ?", 1)).
		Error; err != nil {
		return err
	}

	var count int64

	if err := config.DB.
		Model(&models.PlaylistVideo{}).
		Where("playlist_id = ?", playlistVideo.PlaylistID).
		Count(&count).Error; err != nil {
		return err
	}

	if count == 1 {
		if err := config.DB.
			Model(&models.Playlist{}).
			Where("id = ?", playlistVideo.PlaylistID).
			Update("thumbnail", video.ThumbnailPath).
			Error; err != nil {
			return err
		}
	}

	return nil
}

func (s *playlistVideoService) GetPlaylistDetail(playlistID uint) (*schemas.PlaylistDetailResponse, error) {
	var playlist models.Playlist
	if err := config.DB.
		Where("id = ? AND is_private = ?", playlistID, false).
		First(&playlist).Error; err != nil {
		return nil, ErrPlaylistNotFound
	}

	var videos []models.Video

	err := config.DB.
		Table("playlist_videos").
		Select(`
			videos.id, 
			videos.title, 
			videos.file_path, 
			videos.thumbnail_path,
			videos.duration_video, 
			channels.name, 
			channels.profile_image`,
		).
		Joins("JOIN videos ON videos.id = playlist_videos.video_id").
		Joins("JOIN channels ON channels.id = videos.channel_id").
		Where("playlist_videos.playlist_id = ?", playlistID).
		Scan(&videos).Error

	if err != nil {
		return nil, err
	}

	var response schemas.PlaylistDetailResponse

	response.ID = playlist.ID
	response.Title = playlist.Title
	response.Description = playlist.Description
	response.Thumbnail = playlist.Thumbnail
	response.VideoCount = int64(len(videos))

	for _, v := range videos {
		response.Videos = append(response.Videos, schemas.PlaylistVideoResponse{
			VideoID:       v.ID,
			Title:         v.Title,
			FilePath:      v.FilePath,
			ThumbnailPath: v.ThumbnailPath,
			DurationVideo: v.DurationVideo,
			ChannelName:   v.Channel.Name,
			ProfileImage:  v.Channel.ProfileImage,
		})
	}

	return &response, nil
}

func (s *playlistVideoService) GetMyPlaylistVideo(playlistID, userID uint) (*schemas.MyPlaylistDetailResponse, error) {
	var playlist models.Playlist
	if err := config.DB.
		Where("id = ? AND user_id = ?", playlistID, userID).
		First(&playlist).Error; err != nil {
		return nil, ErrPlaylistNotFound
	}

	var videos []models.Video

	err := config.DB.
		Table("playlist_videos").
		Select(`
		videos.id as video_id,
		videos.title,
		videos.file_path,
		videos.thumbnail_path,
		videos.duration_video,
		videos.views,
		channels.id as channel_id,
		channels.name as channel_name,
		channels.profile_image
	`).
		Joins("JOIN videos ON videos.id = playlist_videos.video_id").
		Joins("JOIN channels ON channels.id = videos.channel_id").
		Where("playlist_videos.playlist_id = ?", playlistID).
		Scan(&videos).Error

	if err != nil {
		return nil, err
	}

	var response schemas.MyPlaylistDetailResponse

	response.ID = playlist.ID
	response.Title = playlist.Title
	response.Description = playlist.Description
	response.Thumbnail = playlist.Thumbnail
	response.IsPrivate = playlist.IsPrivate
	response.VideoCount = int64(len(videos))
	response.CreatedAt = playlist.CreatedAt

	for _, v := range videos {
		response.Videos = append(response.Videos, schemas.MyPlaylistVideoResponse{
			VideoID:       v.ID,
			Title:         v.Title,
			FilePath:      v.FilePath,
			ThumbnailPath: v.ThumbnailPath,
			DurationVideo: v.DurationVideo,
			Views:         v.Views,
			ChannelID:     v.Channel.ID,
			ChannelName:   v.Channel.Name,
			ProfileImage:  v.Channel.ProfileImage,
		})
	}

	return &response, nil
}

func (s *playlistVideoService) DeletePlaylistVideo(ID, userID uint) error {
	var playlistVideo models.PlaylistVideo
	if err := config.DB.First(&playlistVideo, ID).Error; err != nil {
		return ErrPlaylistVideoNotFound
	}

	var playlist models.Playlist
	if err := config.DB.
		Where("id = ? AND user_id = ?", playlistVideo.PlaylistID, userID).
		First(&playlist).Error; err != nil {
		return ErrPlaylistNotFound
	}

	if err := config.DB.Delete(&playlistVideo).Error; err != nil {
		return err
	}

	if err := config.DB.
		Model(&models.Playlist{}).
		Where("id = ?", playlistVideo.PlaylistID).
		UpdateColumn("video_count", gorm.Expr("GREATEST(video_count - 1, 0)")).
		Error; err != nil {
		return err
	}

	var firstVideo models.Video

	err := config.DB.
		Table("playlist_videos").
		Select("videos.thumbnail_path").
		Joins("JOIN videos ON videos.id = playlist_videos.video_id").
		Where("playlist_videos.playlist_id = ?", playlistVideo.PlaylistID).
		First(&firstVideo).Error

	if err != nil {

		config.DB.
			Model(&models.Playlist{}).
			Where("id = ?", playlistVideo.PlaylistID).
			Update("thumbnail", "")

	} else {

		config.DB.
			Model(&models.Playlist{}).
			Where("id = ?", playlistVideo.PlaylistID).
			Update("thumbnail", firstVideo.ThumbnailPath)
	}

	return nil
}
