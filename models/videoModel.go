package models

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	ChannelID     uint     `json:"channel_id"`
	Channel       Channel  `gorm:"foreignKey:ChannelID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	LanguageID    uint     `json:"language_id"`
	Language      Language `gorm:"foreignKey:LanguageID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CategoryID    uint     `json:"category_id"`
	Category      Category `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Title         string   `json:"title" gorm:"size:255;not null"`
	Description   string   `json:"description" gorm:"type:text"`
	FilePath      string   `json:"file_path" gorm:"size:255;not null"`
	ThumbnailPath string   `json:"thumbnail_path" gorm:"size:255;not null"`
	Tags          string   `json:"tags" gorm:"type:text"`
	Resolution    string   `json:"resolution" gorm:"size:50"`
	Size          int64    `json:"size"`
	Views         int64    `json:"views" gorm:"default:0;not null"`
	LikeCount     int64    `json:"like_count" gorm:"default:0;not null"`
	CommentCount  int64    `json:"comment_count" gorm:"default:0"`
	DislikeCount  int64    `json:"dislike_count" gorm:"default:0;not null"`
	DurationVideo int64    `json:"duration_video"`
	Visibility    string   `json:"visibility" gorm:"size:100;default:'public'"`
	Status        string   `json:"status" gorm:"size:20;default:'processing'"`
	ShareCount    int64    `json:"share_count" gorm:"default:0"`
	DownloadCount int64    `json:"download_count" gorm:"default:0"`
}
