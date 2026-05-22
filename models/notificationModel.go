package models

import "gorm.io/gorm"

type Notification struct {
	gorm.Model
	UserID     uint   `json:"user_id"`
	User       User   `json:"-" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	FromUserID *uint  `json:"from_user_id"`
	FromUser   *User  `json:"-" gorm:"foreignKey:FromUserID;constraint:OnUpdate:SET NULL,OnDelete:SET NULL;"`
	Type       string `json:"type" gorm:"size:100"`
	VideoID    *uint  `json:"video_id"`
	CommentID  *uint  `json:"comment_id"`
	ChannelID  *uint  `json:"channel_id"`
	IsRead     bool   `json:"is_read" gorm:"default:false"`
}
