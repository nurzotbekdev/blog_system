package models

import "gorm.io/gorm"

type Subscription struct {
	gorm.Model
	UserID    uint    `json:"user_id" gorm:"index;uniqueIndex:idx_user_channel"`
	User      User    `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ChannelID uint    `json:"channel_id" gorm:"index;uniqueIndex:idx_user_channel"`
	Channel   Channel `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
