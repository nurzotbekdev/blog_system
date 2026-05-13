package models

import "gorm.io/gorm"

type VideoQuality struct {
	gorm.Model
	VideoID  uint   `json:"video_id" gorm:"not null;uniqueIndex:uniq_video_quality"`
	Video    Video  `json:"-" gorm:"foreignKey:VideoID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Quality  string `json:"quality" gorm:"size:20;not null;uniqueIndex:uniq_video_quality"`
	VideoURL string `json:"video_url" gorm:"size:255;not null"`
	Size     int64  `json:"size"`
	Format   string `json:"format" gorm:"size:20"`
}
