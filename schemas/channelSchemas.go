package schemas

import (
	"mime/multipart"
	"time"
)

type CreateChannelRequest struct {
	UserID      uint
	Name        string
	Description string
	ProfileFile *multipart.FileHeader
	BannerFile  *multipart.FileHeader
}

type ChannelResponse struct {
	ID               uint      `json:"id"`
	UserID           uint      `json:"user_id"`
	Email            string    `json:"email"`
	FullName         string    `json:"full_name"`
	AvatarImage      string    `json:"avatar_image"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	ProfileImage     string    `json:"profile_image"`
	BannerImage      string    `json:"banner_image"`
	TotalSubscribers uint      `json:"total_subscribers"`
	TotalVideos      uint      `json:"total_videos"`
	TotalComments    uint      `json:"total_comments"`
	TotalViews       uint      `json:"total_views"`
	TotalWatchTime   float64   `json:"total_watch_time"`
	CreatedAt        time.Time `json:"created_at"`
}

type ChannelSearchResponse struct {
	ID               uint   `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	AvatarImage      string `json:"avatar_image"`
	ProfileImage     string `json:"profile_image"`
	BannerImage      string `json:"banner_image"`
	TotalSubscribers uint   `json:"total_subscribers"`
	TotalVideos      uint   `json:"total_videos"`
	TotalViews       uint   `json:"total_views"`
}

type UpdateChannelRequest struct {
	Name        *string `form:"name"`
	Description *string `form:"description"`
	AvatarImage *string `form:"avatar_image"`
	BannerImage *string `form:"banner_image"`
}
