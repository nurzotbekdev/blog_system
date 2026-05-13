package models

import "gorm.io/gorm"

type VideoDownload struct {
	gorm.Model
	UserID  uint  `json:"user_id"`
	User    User  `json:"-" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	VideoID uint  `json:"video_id"`
	Video   Video `json:"-" gorm:"foreignKey:VideoID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
