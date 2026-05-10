package models

import "gorm.io/gorm"

type Channel struct {
	gorm.Model
	UserID           uint   `json:"user_id" gorm:"index;uniqueIndex"`
	User             User   `json:"-" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name             string `json:"name" gorm:"size:255;unique;not null"`
	Description      string `json:"description" gorm:"type:text"`
	ProfileImage     string `json:"profile_image" gorm:"size:255"`
	BannerImage      string `json:"banner_image" gorm:"size:255"`
	TotalSubscribers uint   `json:"total_subscribers" gorm:"default:0"`
	TotalVideos      uint   `json:"total_videos" gorm:"default:0"`
	TotalComments    uint   `json:"total_comments" gorm:"default:0"`
	TotalViews       uint   `json:"total_views" gorm:"default:0"`
}
