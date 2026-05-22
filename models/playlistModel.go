package models

import "gorm.io/gorm"

type Playlist struct {
	gorm.Model
	UserID      uint   `json:"user_id"`
	User        User   `json:"-" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Title       string `json:"title" gorm:"size:255;not null"`
	Description string `json:"description" gorm:"type:text"`
	IsPrivate   bool   `json:"is_private" gorm:"default:false"`
	Thumbnail   string `json:"thumbnail" gorm:"size:255;default:null"`
	VideoCount  int64  `json:"video_count" gorm:"default:0"`
}
