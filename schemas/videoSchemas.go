package schemas

import (
	"mime/multipart"
	"time"
)

type CreateVideoRequest struct {
	LanguageID    uint                  `form:"language_id"`
	CategoryID    uint                  `form:"category_id"`
	Title         string                `form:"title"`
	Description   string                `form:"description"`
	FilePath      *multipart.FileHeader `form:"file_path"`
	ThumbnailPath *multipart.FileHeader `form:"thumbnail_path"`
}

type VideoQualityResponse struct {
	Quality  string `json:"quality"`
	VideoURL string `json:"video_url"`
	Size     int64  `json:"size"`
	Format   string `json:"format"`
}

type MyVideoResponse struct {
	VideoID       uint                   `json:"video_id"`
	ChannelName   string                 `json:"channel_name"`
	ProfileImage  string                 `json:"profile_image"`
	LanguageName  string                 `json:"language_name"`
	Code          string                 `json:"code"`
	CategoryName  string                 `json:"category_name"`
	Title         string                 `json:"title"`
	Description   string                 `json:"description"`
	FilePath      string                 `json:"file_path"`
	ThumbnailPath string                 `json:"thumbnail_path"`
	Tags          string                 `json:"tags"`
	Resolution    string                 `json:"resolution"`
	Size          int64                  `json:"size"`
	Views         int64                  `json:"views"`
	LikeCount     int64                  `json:"like_count"`
	CommentCount  int64                  `json:"comment_count"`
	DislikeCount  int64                  `json:"dislike_count"`
	DurationVideo int64                  `json:"duration_video"`
	Visibility    string                 `json:"visibility"`
	Status        string                 `json:"status"`
	ShareCount    int64                  `json:"share_count"`
	DownloadCount int64                  `json:"download_count"`
	CreatedAt     time.Time              `json:"created_at"`
	Qualities     []VideoQualityResponse `json:"qualities"`
}

type VideoResponse struct {
	VideoID       uint                   `json:"video_id"`
	ChannelName   string                 `json:"channel_name"`
	ProfileImage  string                 `json:"profile_image"`
	LanguageName  string                 `json:"language_name"`
	Code          string                 `json:"code"`
	CategoryName  string                 `json:"category_name"`
	Title         string                 `json:"title"`
	Description   string                 `json:"description"`
	FilePath      string                 `json:"file_path"`
	ThumbnailPath string                 `json:"thumbnail_path"`
	Tags          string                 `json:"tags"`
	Resolution    string                 `json:"resolution"`
	Size          int64                  `json:"size"`
	Views         int64                  `json:"views"`
	LikeCount     int64                  `json:"like_count"`
	DurationVideo int64                  `json:"duration_video"`
	CreatedAt     time.Time              `json:"created_at"`
	Qualities     []VideoQualityResponse `json:"qualities"`
}

type VideoListResponse struct {
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	Total      int64           `json:"total"`
	TotalPages int             `json:"total_pages"`
	Data       []VideoResponse `json:"data"`
}
