package models

import "gorm.io/gorm"

type History struct {
	gorm.Model
	UserID  uint  `json:"user_id" gorm:"uniqueIndex:uniq_user_video"`
	User    User  `json:"-" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	VideoID uint  `json:"video_id" gorm:"uniqueIndex:uniq_user_video"`
	Video   Video `json:"-" gorm:"foreignKey:VideoID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
