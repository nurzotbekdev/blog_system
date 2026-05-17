package models

import "gorm.io/gorm"

type Like struct {
	gorm.Model
	UserID  uint  `json:"user_id" gorm:"uniqueIndex:uniq_user_video"`
	User    User  `json:"-" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	VideoID uint  `json:"video_id" gorm:"uniqueIndex:uniq_user_video"`
	Video   Video `json:"-" gorm:"foreignKey:VideoID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	IsLike  bool  `json:"is_like" gorm:"not null"`
}
