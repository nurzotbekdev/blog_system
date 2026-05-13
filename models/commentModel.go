package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	UserID       uint      `json:"user_id" gorm:"not null;index"`
	User         User      `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	VideoID      uint      `json:"video_id" gorm:"not null;index"`
	Video        Video     `json:"-" gorm:"foreignKey:VideoID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ParentID     *uint     `json:"parent_id" gorm:"index"`
	Parent       *Comment  `json:"parent,omitempty" gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE;"`
	Content      string    `json:"content" gorm:"type:text;not null"`
	LikeCount    int64     `json:"like_count" gorm:"default:0"`
	DislikeCount int64     `json:"dislike_count" gorm:"default:0"`
	Replies      []Comment `json:"replies,omitempty" gorm:"foreignKey:ParentID"`
}
