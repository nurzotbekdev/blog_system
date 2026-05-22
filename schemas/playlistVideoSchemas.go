package schemas

import "time"

type PlaylistVideoRequest struct {
	PlaylistID uint `json:"playlist_id"`
	VideoID    uint `json:"video_id"`
}

type PlaylistVideoResponse struct {
	VideoID       uint   `json:"video_id"`
	Title         string `json:"title"`
	FilePath      string `json:"file_path"`
	ThumbnailPath string `json:"thumbnail_path"`
	DurationVideo int64  `json:"duration_video"`
	ChannelName   string `json:"channel_name"`
	ProfileImage  string `json:"profile_image"`
}

type PlaylistDetailResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
	VideoCount  int64  `json:"video_count"`

	Videos []PlaylistVideoResponse `json:"videos"`
}

type MyPlaylistVideoResponse struct {
	VideoID       uint   `json:"video_id"`
	Title         string `json:"title"`
	FilePath      string `json:"file_path"`
	ThumbnailPath string `json:"thumbnail_path"`
	DurationVideo int64  `json:"duration_video"`
	Views         int64  `json:"views"`
	ChannelID     uint   `json:"channel_id"`
	ChannelName   string `json:"channel_name"`
	ProfileImage  string `json:"profile_image"`
}

type MyPlaylistDetailResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Thumbnail   string    `json:"thumbnail"`
	IsPrivate   bool      `json:"is_private"`
	VideoCount  int64     `json:"video_count"`
	CreatedAt   time.Time `json:"created_at"`

	Videos []MyPlaylistVideoResponse `json:"videos"`
}
