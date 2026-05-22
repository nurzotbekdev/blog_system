package schemas

import "time"

type NotificationPayload struct {
	UserID     uint   `json:"user_id"`
	FromUserID *uint  `json:"from_user_id"`
	Type       string `json:"type"`
	VideoID    *uint  `json:"video_id"`
	CommentID  *uint  `json:"comment_id"`
	ChannelID  *uint  `json:"channel_id"`
}

type NotificationResponse struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	FromUserID *uint     `json:"from_user_id"`
	VideoID    *uint     `json:"video_id"`
	CommentID  *uint     `json:"comment_id"`
	ChannelID  *uint     `json:"channel_id"`
	CreatedAt  time.Time `json:"created_at"`
}
