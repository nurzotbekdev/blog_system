package models

import "gorm.io/gorm"

type CommentReaction struct {
	gorm.Model
	UserID    uint    `json:"user_id" gorm:"uniqueIndex:uniq_user_comment"`
	User      User    `json:"-" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CommentID uint    `json:"comment_id" gorm:"uniqueIndex:uniq_user_comment"`
	Comment   Comment `json:"-" gorm:"foreignKey:CommentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	IsLike    bool    `json:"is_like" gorm:"not null"`
}
