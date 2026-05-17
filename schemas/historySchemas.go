package schemas

import "time"

type HistoryResponse struct {
	ID            uint      `json:"id"`
	VideoID       uint      `json:"video_id"`
	Title         string    `json:"title"`
	ThumbnailPath string    `json:"thumbnail_path"`
	DurationVideo int64     `json:"duration_video"`
	ChannelName   string    `json:"channel_name"`
	ProfileImage  string    `json:"profile_image"`
	CreatedAt     time.Time `json:"created_at"`
}
